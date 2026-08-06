package domain

// =============================================================================
// Calculator Domain Models — Exact parity with Java echallan-calculator
// =============================================================================

// RequestInfo — same as echallan-services (shared DIGIT contract)
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

type Citizen struct {
	Id           int    `json:"id,omitempty"`
	Uuid         string `json:"uuid,omitempty"`
	UserName     string `json:"userName,omitempty"`
	Name         string `json:"name,omitempty"`
	MobileNumber string `json:"mobileNumber,omitempty"`
	Type         string `json:"type,omitempty"`
	TenantId     string `json:"tenantId,omitempty"`
	Roles        []Role `json:"roles,omitempty"`
}

type Role struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	TenantId string `json:"tenantId,omitempty"`
}

type ResponseInfo struct {
	ApiId    string `json:"apiId"`
	Ver      string `json:"ver"`
	Ts       *int64 `json:"ts"`
	ResMsgId string `json:"resMsgId,omitempty"`
	MsgId    string `json:"msgId"`
	Status   string `json:"status"`
}

type Error struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
	Params      string `json:"params,omitempty"`
}

type ErrorResponse struct {
	ResponseInfo *ResponseInfo `json:"ResponseInfo"`
	Errors       []Error       `json:"Errors"`
}

// =============================================================================
// Calculator-specific models
// =============================================================================

// Challan — simplified version used in calculation criteria
type Challan struct {
	TenantId          string   `json:"tenantId"`
	BusinessService   string   `json:"businessService"`
	ChallanNo         string   `json:"challanNo,omitempty"`
	ApplicationStatus string   `json:"applicationStatus,omitempty"`
	Description       string   `json:"description,omitempty"`
	Citizen           *Citizen `json:"citizen,omitempty"`
	Amount            []Amount `json:"amount"`
	TaxPeriodFrom     int64    `json:"taxPeriodFrom"`
	TaxPeriodTo       int64    `json:"taxPeriodTo"`
}

type Amount struct {
	TaxHeadCode string  `json:"taxHeadCode"`
	Amount      float64 `json:"amount"`
}

// CalculationRequest is the inbound JSON body for /_calculate.
// CRITICAL: The JSON key is "CalulationCriteria" (INTENTIONAL TYPO — missing 'c')
// This typo exists in the original Java code and MUST be preserved for backward compatibility.
// Verified in E2E report Step 9.
type CalculationRequest struct {
	RequestInfo        *RequestInfo          `json:"RequestInfo"`
	CalulationCriteria []CalulationCriterion `json:"CalulationCriteria"`
}

// CalulationCriterion — note the intentional typo in the name to match Java.
type CalulationCriterion struct {
	Challan  *Challan `json:"challan"`
	TenantId string   `json:"tenantId"`
}

// CalculationResponse is the outbound JSON for /_calculate.
type CalculationResponse struct {
	ResponseInfo *ResponseInfo `json:"ResponseInfo"`
	Calculations []Calculation `json:"Calculations,omitempty"`
	Errors       []Error       `json:"Errors,omitempty"`
}

type Calculation struct {
	TenantId        string         `json:"tenantId"`
	ChallanNo       string         `json:"challanNo"`
	TaxHeadEstimate []TaxEstimate  `json:"taxHeadEstimate,omitempty"`
}

type TaxEstimate struct {
	TaxHeadCode    string  `json:"taxHeadCode"`
	EstimateAmount float64 `json:"estimateAmount"`
}

// =============================================================================
// Demand models — for calling billing-service
// =============================================================================

type DemandRequest struct {
	RequestInfo *RequestInfo `json:"RequestInfo"`
	Demands     []Demand     `json:"Demands"`
}

type Demand struct {
	TenantId          string         `json:"tenantId"`
	BusinessService   string         `json:"businessService"`
	ConsumerCode      string         `json:"consumerCode"`
	ConsumerType      string         `json:"consumerType"`
	TaxPeriodFrom     int64          `json:"taxPeriodFrom"`
	TaxPeriodTo       int64          `json:"taxPeriodTo"`
	DemandDetails     []DemandDetail `json:"demandDetails"`
	Payer             *Citizen       `json:"payer,omitempty"`
}

type DemandDetail struct {
	TaxHeadMasterCode string  `json:"taxHeadMasterCode"`
	TaxAmount         float64 `json:"taxAmount"`
	CollectionAmount  float64 `json:"collectionAmount"`
}

type DemandResponse struct {
	ResponseInfo *ResponseInfo `json:"ResponseInfo"`
	Demands      []Demand      `json:"Demands"`
}
