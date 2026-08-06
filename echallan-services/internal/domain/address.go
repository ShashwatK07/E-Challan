package domain

type Address struct {
	Id           string    `json:"id,omitempty"`
	TenantId     string    `json:"tenantId,omitempty"`
	DoorNo       string    `json:"doorNo,omitempty"`
	PlotNo       string    `json:"plotNo,omitempty"`
	Landmark     string    `json:"landmark,omitempty"`
	City         string    `json:"city,omitempty"`
	Pincode      string    `json:"pincode,omitempty"`
	Detail       string    `json:"detail,omitempty"`
	BuildingName string    `json:"buildingName,omitempty"`
	Street       string    `json:"street,omitempty"`
	Latitude     float64   `json:"latitude,omitempty"`
	Longitude    float64   `json:"longitude,omitempty"`
	Locality     *Boundary `json:"locality,omitempty"`
}

type Boundary struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

type ChallanAddressDB struct {
	Id           string  `gorm:"column:id;primaryKey"`
	EchallanId   string  `gorm:"column:echallanid"`
	DoorNo       string  `gorm:"column:doorno"`
	BuildingName string  `gorm:"column:buildingname"`
	Street       string  `gorm:"column:street"`
	City         string  `gorm:"column:city"`
	Pincode      string  `gorm:"column:pincode"`
	Latitude     float64 `gorm:"column:latitude"`
	Longitude    float64 `gorm:"column:longitude"`
	LocalityCode string  `gorm:"column:locality_code"`
	TenantId     string  `gorm:"column:tenantid"`
}

func (ChallanAddressDB) TableName() string { return "eg_challan_address" }

func (db *ChallanAddressDB) ToAddress() *Address {
	return &Address{
		Id:           db.Id,
		TenantId:     db.TenantId,
		DoorNo:       db.DoorNo,
		BuildingName: db.BuildingName,
		Street:       db.Street,
		City:         db.City,
		Pincode:      db.Pincode,
		Latitude:     db.Latitude,
		Longitude:    db.Longitude,
		Locality:     &Boundary{Code: db.LocalityCode},
	}
}

func (a *Address) ToAddressDB(echallanId string) *ChallanAddressDB {
	db := &ChallanAddressDB{
		Id:           a.Id,
		EchallanId:   echallanId,
		DoorNo:       a.DoorNo,
		BuildingName: a.BuildingName,
		Street:       a.Street,
		City:         a.City,
		Pincode:      a.Pincode,
		Latitude:     a.Latitude,
		Longitude:    a.Longitude,
		TenantId:     a.TenantId,
	}
	if a.Locality != nil {
		db.LocalityCode = a.Locality.Code
	}
	return db
}
