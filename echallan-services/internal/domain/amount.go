package domain

type Amount struct {
	TaxHeadCode string  `json:"taxHeadCode"`
	Amount      float64 `json:"amount"`
}
