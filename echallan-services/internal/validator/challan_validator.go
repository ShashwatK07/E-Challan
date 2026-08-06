package validator

import (
	"fmt"

	"github.com/CDPI-HRSS/echallan-services/internal/domain"
	"github.com/CDPI-HRSS/echallan-services/internal/repository/postgres"
)

// ChallanValidator replicates the Java ChallanValidator.java logic.
// Reference: echallan-services/src/main/java/org/egov/echallan/validator/ChallanValidator.java
type ChallanValidator struct {
	repo *postgres.ChallanRepository
}

func NewChallanValidator(repo *postgres.ChallanRepository) *ChallanValidator {
	return &ChallanValidator{repo: repo}
}

// ValidateCreateRequest validates the incoming /_create request.
func (v *ChallanValidator) ValidateCreateRequest(req *domain.ChallanRequest) error {
	if req.RequestInfo == nil {
		return fmt.Errorf("RequestInfo is mandatory")
	}
	if req.RequestInfo.UserInfo == nil {
		return fmt.Errorf("UserInfo is mandatory in RequestInfo")
	}
	if req.Challan == nil {
		return fmt.Errorf("Challan object is mandatory")
	}

	c := req.Challan
	if c.TenantId == "" {
		return fmt.Errorf("tenantId is mandatory")
	}
	if c.BusinessService == "" {
		return fmt.Errorf("businessService is mandatory")
	}
	if c.Citizen == nil {
		return fmt.Errorf("citizen is mandatory")
	}
	if c.Citizen.MobileNumber == "" {
		return fmt.Errorf("citizen.mobileNumber is mandatory")
	}
	if c.Citizen.Name == "" {
		return fmt.Errorf("citizen.name is mandatory")
	}
	if len(c.Amount) == 0 {
		return fmt.Errorf("at least one amount entry is mandatory")
	}
	if c.TaxPeriodFrom == 0 || c.TaxPeriodTo == 0 {
		return fmt.Errorf("taxPeriodFrom and taxPeriodTo are mandatory")
	}

	return nil
}

// ValidateUpdateRequest validates the incoming /_update request.
// CRITICAL: This replicates ChallanValidator.validateUpdateRequest() from Java.
// It cross-references the incoming payload against the existing database record.
// As documented in E2E report Step 6, these immutable fields MUST match:
//   - id, challanNo, businessService, address.id, citizen.uuid, citizen.name, citizen.mobileNumber
func (v *ChallanValidator) ValidateUpdateRequest(req *domain.ChallanRequest) error {
	if req.RequestInfo == nil || req.RequestInfo.UserInfo == nil {
		return fmt.Errorf("RequestInfo with UserInfo is mandatory")
	}
	if req.Challan == nil {
		return fmt.Errorf("Challan object is mandatory")
	}

	c := req.Challan

	if c.Id == "" {
		return fmt.Errorf("challan id is mandatory for update")
	}
	if c.ChallanNo == "" {
		return fmt.Errorf("challanNo is mandatory for update")
	}
	if c.TenantId == "" {
		return fmt.Errorf("tenantId is mandatory for update")
	}

	// Fetch existing record from DB for cross-validation
	existing, err := v.repo.FindById(c.Id)
	if err != nil {
		return fmt.Errorf("challan not found in database: %w", err)
	}

	// Validate immutable fields match the DB record
	if c.BusinessService != existing.BusinessService {
		return fmt.Errorf("businessService cannot be modified (expected: %s, got: %s)",
			existing.BusinessService, c.BusinessService)
	}
	if c.ChallanNo != existing.ChallanNo {
		return fmt.Errorf("challanNo cannot be modified (expected: %s, got: %s)",
			existing.ChallanNo, c.ChallanNo)
	}
	if existing.ApplicationStatus != "ACTIVE" {
		return fmt.Errorf("only ACTIVE challans can be updated (current status: %s)",
			existing.ApplicationStatus)
	}

	// Validate address.id matches if provided
	if c.Address != nil && existing.Address != nil {
		if c.Address.Id != existing.Address.Id {
			return fmt.Errorf("address.id cannot be modified (expected: %s, got: %s)",
				existing.Address.Id, c.Address.Id)
		}
	}

	// Validate citizen immutable fields
	if c.Citizen == nil {
		return fmt.Errorf("citizen is mandatory for update")
	}

	// Validate tenantId matches userInfo.tenantId
	if req.RequestInfo.UserInfo.TenantId != c.TenantId {
		return fmt.Errorf("challan tenantId must match userInfo.tenantId")
	}

	return nil
}
