package postgres

import (
	"fmt"
	"log"

	"github.com/CDPI-HRSS/echallan-services/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ChallanRepository handles all database operations for challans.
// This replaces the Java ChallanRowMapper + ChallanQueryBuilder pattern.
type ChallanRepository struct {
	db *gorm.DB
}

// NewChallanRepository creates a new repository with a GORM database connection.
func NewChallanRepository(host, port, user, password, dbname string) (*ChallanRepository, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	log.Printf("✅ Connected to PostgreSQL: %s@%s:%s/%s", user, host, port, dbname)
	return &ChallanRepository{db: db}, nil
}

// Save inserts a new challan + its address into the database.
// This replicates the Java "dual persistence" pattern — direct JDBC INSERT.
func (r *ChallanRepository) Save(challan *domain.Challan) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Insert into eg_echallan
		challanDB := challan.ToChallanDB()
		if err := tx.Create(challanDB).Error; err != nil {
			return fmt.Errorf("failed to insert challan: %w", err)
		}

		// 2. Insert into eg_challan_address
		if challan.Address != nil {
			addressDB := challan.Address.ToAddressDB(challan.Id)
			if err := tx.Create(addressDB).Error; err != nil {
				return fmt.Errorf("failed to insert address: %w", err)
			}
		}

		return nil
	})
}

// Update modifies an existing challan record.
func (r *ChallanRepository) Update(challan *domain.Challan) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		challanDB := challan.ToChallanDB()
		result := tx.Model(&domain.ChallanDB{}).
			Where("id = ?", challan.Id).
			Updates(map[string]interface{}{
				"description":       challanDB.Description,
				"applicationstatus": challanDB.ApplicationStatus,
				"lastmodifiedby":    challanDB.LastModifiedBy,
				"lastmodifiedtime":  challanDB.LastModifiedTime,
			})
		if result.Error != nil {
			return fmt.Errorf("failed to update challan: %w", result.Error)
		}
		return nil
	})
}

// Search retrieves challans matching the given criteria.
// This replicates the Java ChallanQueryBuilder.getChallanSearchQuery() method.
func (r *ChallanRepository) Search(criteria domain.SearchCriteria) ([]*domain.Challan, error) {
	query := r.db.Model(&domain.ChallanDB{})

	if criteria.TenantId != "" {
		query = query.Where("tenantid = ?", criteria.TenantId)
	}
	if criteria.ChallanNo != "" {
		query = query.Where("challanno = ?", criteria.ChallanNo)
	}
	if criteria.AccountId != "" {
		query = query.Where("accountid = ?", criteria.AccountId)
	}
	if criteria.BusinessService != "" {
		query = query.Where("businessservice = ?", criteria.BusinessService)
	}
	if criteria.Status != "" {
		query = query.Where("applicationstatus = ?", criteria.Status)
	}

	// Default limit to prevent unbounded queries
	limit := criteria.Limit
	if limit <= 0 {
		limit = 100
	}
	query = query.Limit(limit).Offset(criteria.Offset)

	var challanDBs []domain.ChallanDB
	if err := query.Find(&challanDBs).Error; err != nil {
		return nil, fmt.Errorf("failed to search challans: %w", err)
	}

	// Convert DB models to API models and enrich with address
	challans := make([]*domain.Challan, 0, len(challanDBs))
	for _, cdb := range challanDBs {
		challan := cdb.ToChallan()

		// Fetch associated address
		var addressDB domain.ChallanAddressDB
		if err := r.db.Where("echallanid = ?", cdb.Id).First(&addressDB).Error; err == nil {
			challan.Address = addressDB.ToAddress()
		}

		challans = append(challans, challan)
	}

	return challans, nil
}

// FindById retrieves a single challan by its UUID.
func (r *ChallanRepository) FindById(id string) (*domain.Challan, error) {
	var challanDB domain.ChallanDB
	if err := r.db.Where("id = ?", id).First(&challanDB).Error; err != nil {
		return nil, fmt.Errorf("challan not found with id: %s", id)
	}

	challan := challanDB.ToChallan()

	// Fetch associated address
	var addressDB domain.ChallanAddressDB
	if err := r.db.Where("echallanid = ?", id).First(&addressDB).Error; err == nil {
		challan.Address = addressDB.ToAddress()
	}

	return challan, nil
}

// FindByChallanNo retrieves a single challan by its challanNo.
func (r *ChallanRepository) FindByChallanNo(challanNo string) (*domain.Challan, error) {
	var challanDB domain.ChallanDB
	if err := r.db.Where("challanno = ?", challanNo).First(&challanDB).Error; err != nil {
		return nil, fmt.Errorf("challan not found with challanNo: %s", challanNo)
	}

	challan := challanDB.ToChallan()

	var addressDB domain.ChallanAddressDB
	if err := r.db.Where("echallanid = ?", challanDB.Id).First(&addressDB).Error; err == nil {
		challan.Address = addressDB.ToAddress()
	}

	return challan, nil
}

// Count returns aggregated challan counts grouped by applicationStatus.
// This replicates the Java ChallanQueryBuilder.getChallanCountQuery() method.
// Returns counts as STRINGS to match Java JSON output exactly (E2E report Step 7).
func (r *ChallanRepository) Count(tenantId string) (*domain.ChallanCount, error) {
	type StatusCount struct {
		Status string
		Count  int64
	}

	var results []StatusCount
	err := r.db.Model(&domain.ChallanDB{}).
		Select("applicationstatus as status, count(*) as count").
		Where("tenantid = ?", tenantId).
		Group("applicationstatus").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count challans: %w", err)
	}

	count := &domain.ChallanCount{
		PaidChallan:      "0",
		CancelledChallan: "0",
		TotalChallan:     "0",
		ActiveChallan:    "0",
	}

	var total int64
	for _, r := range results {
		total += r.Count
		switch r.Status {
		case "ACTIVE":
			count.ActiveChallan = fmt.Sprintf("%d", r.Count)
		case "PAID":
			count.PaidChallan = fmt.Sprintf("%d", r.Count)
		case "CANCELLED":
			count.CancelledChallan = fmt.Sprintf("%d", r.Count)
		}
	}
	count.TotalChallan = fmt.Sprintf("%d", total)

	return count, nil
}
