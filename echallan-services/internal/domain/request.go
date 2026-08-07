package domain

type RequestInfo struct {
	ApiId     string   `json:"apiId"`
	Ver       string   `json:"ver"`
	Ts        *int64   `json:"ts"`
	Action    string   `json:"action,omitempty"`
	Did       string   `json:"did,omitempty"`
	Key       string   `json:"key,omitempty"`
	MsgId     string   `json:"msgId"`
	AuthToken string   `json:"authToken,omitempty"`
	UserInfo  *Citizen `json:"userInfo,omitempty"`
}

type ChallanRequest struct {
	RequestInfo *RequestInfo `json:"RequestInfo"`
	Challan     *Challan     `json:"Challan"`
}

type RequestInfoWrapper struct {
	RequestInfo *RequestInfo `json:"RequestInfo"`
}

type SearchCriteria struct {
	TenantId        string `form:"tenantId"`
	ChallanNo       string `form:"challanNo"`
	AccountId       string `form:"accountId"`
	MobileNumber    string `form:"mobileNumber"`
	BusinessService string `form:"businessService"`
	Status          string `form:"status"`
	Offset          int    `form:"offset"`
	Limit           int    `form:"limit"`
}
