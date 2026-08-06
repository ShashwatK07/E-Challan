package domain

type Calculation struct {
	TenantId         string             `json:"tenantId"`
	TaxHeadEstimates []TaxHeadEstimate `json:"taxHeadEstimates"`
}

type TaxHeadEstimate struct {
	TaxHeadCode    string  `json:"taxHeadCode"`
	EstimateAmount float64 `json:"estimateAmount"`
}
