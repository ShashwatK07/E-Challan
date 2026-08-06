package domain

type Challan struct {
	Citizen           *Citizen      `json:"citizen,omitempty"`
	Id                string        `json:"id,omitempty"`
	TenantId          string        `json:"tenantId"`
	BusinessService   string        `json:"businessService"`
	ChallanNo         string        `json:"challanNo,omitempty"`
	ReferenceId       string        `json:"referenceId,omitempty"`
	Description       string        `json:"description,omitempty"`
	AccountId         string        `json:"accountId,omitempty"`
	AdditionalDetail  interface{}   `json:"additionalDetail,omitempty"`
	Source            string        `json:"source,omitempty"`
	TaxPeriodFrom     int64         `json:"taxPeriodFrom"`
	TaxPeriodTo       int64         `json:"taxPeriodTo"`
	Calculation       interface{}   `json:"calculation,omitempty"`
	Amount            []Amount      `json:"amount"`
	Address           *Address      `json:"address,omitempty"`
	AuditDetails      *AuditDetails `json:"auditDetails,omitempty"`
	ApplicationStatus string        `json:"applicationStatus,omitempty"`
}

type ChallanDB struct {
	Id                string `gorm:"column:id;primaryKey"`
	TenantId          string `gorm:"column:tenantid"`
	BusinessService   string `gorm:"column:businessservice"`
	ChallanNo         string `gorm:"column:challanno"`
	ReferenceId       string `gorm:"column:referenceid"`
	Description       string `gorm:"column:description"`
	AccountId         string `gorm:"column:accountid"`
	Source            string `gorm:"column:source"`
	TaxPeriodFrom     int64  `gorm:"column:taxperiodfrom"`
	TaxPeriodTo       int64  `gorm:"column:taxperiodto"`
	ApplicationStatus string `gorm:"column:applicationstatus"`
	Filestoreid       string `gorm:"column:filestoreid"`
	ReceiptNumber     string `gorm:"column:receiptnumber"`
	CreatedBy         string `gorm:"column:createdby"`
	LastModifiedBy    string `gorm:"column:lastmodifiedby"`
	CreatedTime       int64  `gorm:"column:createdtime"`
	LastModifiedTime  int64  `gorm:"column:lastmodifiedtime"`
}

func (ChallanDB) TableName() string { return "eg_echallan" }

func (db *ChallanDB) ToChallan() *Challan {
	return &Challan{
		Id:                db.Id,
		TenantId:          db.TenantId,
		BusinessService:   db.BusinessService,
		ChallanNo:         db.ChallanNo,
		ReferenceId:       db.ReferenceId,
		Description:       db.Description,
		AccountId:         db.AccountId,
		Source:            db.Source,
		TaxPeriodFrom:     db.TaxPeriodFrom,
		TaxPeriodTo:       db.TaxPeriodTo,
		ApplicationStatus: db.ApplicationStatus,
		AuditDetails: &AuditDetails{
			CreatedBy:        db.CreatedBy,
			LastModifiedBy:   db.LastModifiedBy,
			CreatedTime:      db.CreatedTime,
			LastModifiedTime: db.LastModifiedTime,
		},
	}
}

func (c *Challan) ToChallanDB() *ChallanDB {
	db := &ChallanDB{
		Id:                c.Id,
		TenantId:          c.TenantId,
		BusinessService:   c.BusinessService,
		ChallanNo:         c.ChallanNo,
		ReferenceId:       c.ReferenceId,
		Description:       c.Description,
		AccountId:         c.AccountId,
		Source:            c.Source,
		TaxPeriodFrom:     c.TaxPeriodFrom,
		TaxPeriodTo:       c.TaxPeriodTo,
		ApplicationStatus: c.ApplicationStatus,
	}
	if c.AuditDetails != nil {
		db.CreatedBy = c.AuditDetails.CreatedBy
		db.LastModifiedBy = c.AuditDetails.LastModifiedBy
		db.CreatedTime = c.AuditDetails.CreatedTime
		db.LastModifiedTime = c.AuditDetails.LastModifiedTime
	}
	return db
}

type ChallanCount struct {
	PaidChallan      string `json:"paidChallan"`
	CancelledChallan string `json:"cancelledChallan"`
	TotalChallan     string `json:"totalChallan"`
	ActiveChallan    string `json:"activeChallan"`
}
