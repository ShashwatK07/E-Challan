package domain

type ResponseInfo struct {
	ApiId    string `json:"apiId"`
	Ver      string `json:"ver"`
	Ts       *int64 `json:"ts"`
	ResMsgId string `json:"resMsgId,omitempty"`
	MsgId    string `json:"msgId"`
	Status   string `json:"status"`
}

type ChallanResponse struct {
	CountOfServices      int           `json:"countOfServices"`
	TotalAmountCollected int           `json:"totalAmountCollected"`
	ChallanValidity      int           `json:"challanValidity"`
	ResponseInfo         *ResponseInfo `json:"responseInfo"`
	Challans             []*Challan    `json:"challans"`
	TotalCount           int           `json:"totalCount"`
}

type CountResponse struct {
	ResponseInfo *ResponseInfo `json:"ResponseInfo"`
	ChallanCount *ChallanCount `json:"ChallanCount"`
}

type ErrorResponse struct {
	ResponseInfo *ResponseInfo `json:"ResponseInfo,omitempty"`
	Errors       []Error       `json:"Errors"`
}

type Error struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
	Params      string `json:"params,omitempty"`
}
