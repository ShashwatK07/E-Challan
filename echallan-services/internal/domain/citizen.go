package domain

type Citizen struct {
	Id               int    `json:"id,omitempty"`
	Uuid             string `json:"uuid,omitempty"`
	UserName         string `json:"userName,omitempty"`
	Name             string `json:"name,omitempty"`
	MobileNumber     string `json:"mobileNumber,omitempty"`
	Active           *bool  `json:"active,omitempty"`
	Type             string `json:"type,omitempty"`
	TenantId         string `json:"tenantId,omitempty"`
	Roles            []Role `json:"roles,omitempty"`
	CreatedBy        string `json:"createdBy,omitempty"`
	CreatedDate      int64  `json:"createdDate,omitempty"`
	LastModifiedBy   string `json:"lastModifiedBy,omitempty"`
	LastModifiedDate int64  `json:"lastModifiedDate,omitempty"`
}

type Role struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	TenantId string `json:"tenantId,omitempty"`
}
