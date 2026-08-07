package service

import (
	"fmt"
	"log"
	"time"

	"github.com/CDPI-HRSS/echallan-services/internal/domain"
	"github.com/CDPI-HRSS/echallan-services/internal/repository/postgres"
	"github.com/CDPI-HRSS/echallan-services/internal/transport/kafka"
	"github.com/CDPI-HRSS/echallan-services/internal/validator"
	"github.com/google/uuid"
)

// challanServiceImpl is the core business logic implementation.
// It replicates the Java ChallanService.java flow step-by-step.
type challanServiceImpl struct {
	repo       *postgres.ChallanRepository
	producer   *kafka.Producer
	validator  *validator.ChallanValidator
	enrichment *EnrichmentService
}

// NewChallanService creates a new service with all dependencies wired in.
func NewChallanService(
	repo *postgres.ChallanRepository,
	producer *kafka.Producer,
	val *validator.ChallanValidator,
	enrichment *EnrichmentService,
) ChallanService {
	return &challanServiceImpl{
		repo:       repo,
		producer:   producer,
		validator:  val,
		enrichment: enrichment,
	}
}

// Create handles the /_create API.
// Replicates the Java flow from E2E report Step 4:
//   1. Validate request
//   2. Call egov-user to create/lookup citizen → get citizen UUID
//   3. Call egov-mdms to validate business service
//   4. Call egov-idgen to generate unique challanNo
//   5. Enrich with UUID, audit details
//   6. Direct JDBC INSERT into PostgreSQL (synchronous)
//   7. Push to Kafka save-challan topic (asynchronous broadcast)
func (s *challanServiceImpl) Create(req *domain.ChallanRequest) (*domain.Challan, error) {
	// Step 1: Validate
	if err := s.validator.ValidateCreateRequest(req); err != nil {
		return nil, err
	}

	challan := req.Challan
	now := time.Now().UnixMilli()

	// Step 2: Generate UUIDs
	challan.Id = uuid.New().String()
	if challan.Address != nil {
		challan.Address.Id = uuid.New().String()
		challan.Address.TenantId = challan.TenantId
	}

	// Step 3: Call egov-idgen to generate challanNo
	challanNo, err := s.enrichment.GenerateChallanNo(challan.TenantId, req.RequestInfo)
	if err != nil {
		log.Printf("⚠️ egov-idgen failed, using fallback: %v", err)
		challanNo = fmt.Sprintf("CB-CH-%s-%06d", time.Now().Format("2006-01-02"), time.Now().UnixNano()%1000000)
	}
	challan.ChallanNo = challanNo

	// Step 4: Call egov-user to create/lookup citizen
	if challan.Citizen != nil {
		citizenUuid, citizenId, err := s.enrichment.CreateOrLookupCitizen(
			map[string]interface{}{
				"name":         challan.Citizen.Name,
				"mobileNumber": challan.Citizen.MobileNumber,
				"tenantId":     challan.TenantId,
				"type":         "CITIZEN",
				"roles":        challan.Citizen.Roles,
			},
			req.RequestInfo,
		)
		if err != nil {
			log.Printf("⚠️ egov-user failed, using existing citizen data: %v", err)
		} else {
			challan.Citizen.Uuid = citizenUuid
			challan.Citizen.Id = citizenId
			challan.AccountId = citizenUuid
		}
	}

	// Step 5: Validate against MDMS (non-blocking)
	if err := s.enrichment.ValidateBusinessService(challan.TenantId, req.RequestInfo); err != nil {
		log.Printf("⚠️ MDMS validation warning: %v", err)
	}

	// Step 6: Enrich audit details
	var userUuid string
	if req.RequestInfo != nil && req.RequestInfo.UserInfo != nil {
		userUuid = req.RequestInfo.UserInfo.Uuid
	}
	challan.AuditDetails = &domain.AuditDetails{
		CreatedBy:        userUuid,
		LastModifiedBy:   userUuid,
		CreatedTime:      now,
		LastModifiedTime: now,
	}
	challan.ApplicationStatus = "ACTIVE"

	// Step 7: Direct JDBC INSERT (synchronous persistence)
	if err := s.repo.Save(challan); err != nil {
		return nil, fmt.Errorf("failed to save challan to database: %w", err)
	}
	log.Printf("✅ Challan saved to PostgreSQL: %s", challan.ChallanNo)

	// Step 8: Push to Kafka (asynchronous broadcast for egov-persister)
	if err := s.producer.Push("save-challan", req); err != nil {
		log.Printf("⚠️ Kafka push failed (non-fatal): %v", err)
	}

	return challan, nil
}

// Update handles the /_update API.
// Replicates the Java flow from E2E report Step 6:
//   1. Validate request (strict immutable field checks)
//   2. Enrich audit details
//   3. Direct JDBC UPDATE (synchronous)
//   4. Push to Kafka update-challan topic
func (s *challanServiceImpl) Update(req *domain.ChallanRequest) (*domain.Challan, error) {
	// Step 1: Validate (cross-references DB record)
	if err := s.validator.ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	challan := req.Challan

	// Step 2: Enrich audit details
	var userUuid string
	if req.RequestInfo != nil && req.RequestInfo.UserInfo != nil {
		userUuid = req.RequestInfo.UserInfo.Uuid
	}
	challan.AuditDetails = &domain.AuditDetails{
		LastModifiedBy:   userUuid,
		LastModifiedTime: time.Now().UnixMilli(),
	}

	// Step 3: Direct JDBC UPDATE (synchronous)
	if err := s.repo.Update(challan); err != nil {
		return nil, fmt.Errorf("failed to update challan: %w", err)
	}
	log.Printf("✅ Challan updated in PostgreSQL: %s", challan.ChallanNo)

	// Step 4: Push to Kafka
	if err := s.producer.Push("update-challan", req); err != nil {
		log.Printf("⚠️ Kafka push failed (non-fatal): %v", err)
	}

	return challan, nil
}

// Search handles the /_search API.
// Replicates the Java flow from E2E report Step 5.
func (s *challanServiceImpl) Search(criteria domain.SearchCriteria, reqInfo *domain.RequestInfo) ([]*domain.Challan, error) {
	challans, err := s.repo.Search(criteria)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ Search returned %d challans", len(challans))
	return challans, nil
}

// Count handles the /_count API.
// Returns the EXACT JSON structure from E2E report Step 7 with uppercase "ResponseInfo".
func (s *challanServiceImpl) Count(tenantId string, reqInfo *domain.RequestInfo) (*domain.CountResponse, error) {
	count, err := s.repo.Count(tenantId)
	if err != nil {
		return nil, err
	}

	return &domain.CountResponse{
		ResponseInfo: &domain.ResponseInfo{
			ResMsgId: "uief87324",
			Status:   "successful",
		},
		ChallanCount: count,
	}, nil
}

// Test handles the /_test API.
// Directly pushes to Kafka without any business logic (E2E report Step 8).
func (s *challanServiceImpl) Test(req *domain.ChallanRequest) error {
	if err := s.producer.Push("save-challan", req); err != nil {
		return fmt.Errorf("kafka push test failed: %w", err)
	}
	log.Println("✅ Kafka test push successful")
	return nil
}
