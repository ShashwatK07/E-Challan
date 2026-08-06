package service

import "github.com/CDPI-HRSS/echallan-services/internal/domain"

// ChallanService defines the interface for all challan business operations.
// This is the Go equivalent of the Java ChallanService interface.
type ChallanService interface {
	Create(req *domain.ChallanRequest) (*domain.Challan, error)
	Update(req *domain.ChallanRequest) (*domain.Challan, error)
	Search(criteria domain.SearchCriteria, reqInfo *domain.RequestInfo) ([]*domain.Challan, error)
	Count(tenantId string, reqInfo *domain.RequestInfo) (*domain.CountResponse, error)
	Test(req *domain.ChallanRequest) error
}
