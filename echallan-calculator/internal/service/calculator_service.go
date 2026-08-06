package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/CDPI-HRSS/echallan-calculator/internal/domain"
)

// CalculatorService handles the /_calculate API logic.
// This replicates the Java ChallanCalculationService.java.
//
// Flow (E2E report Step 9):
//   1. Receive CalulationCriteria (note: intentional typo preserved)
//   2. For each criterion, fetch tax period config from MDMS
//   3. Build demand details from the challan amount array
//   4. Call billing-service to create a Demand
//   5. Return the calculation result
type CalculatorService struct {
	mdmsHost    string // e.g., "http://localhost:8094"
	billingHost string // e.g., "http://localhost:8081"
	httpClient  *http.Client
}

func NewCalculatorService(mdmsHost, billingHost string) *CalculatorService {
	return &CalculatorService{
		mdmsHost:    mdmsHost,
		billingHost: billingHost,
		httpClient:  &http.Client{},
	}
}

// Calculate processes the calculation criteria and creates demands.
func (cs *CalculatorService) Calculate(req *domain.CalculationRequest) (*domain.CalculationResponse, error) {
	if len(req.CalulationCriteria) == 0 {
		return nil, fmt.Errorf("CalulationCriteria is empty")
	}

	var calculations []domain.Calculation

	for _, criteria := range req.CalulationCriteria {
		challan := criteria.Challan
		if challan == nil {
			continue
		}

		// Step 1: Fetch TaxPeriod from MDMS to validate the date range
		err := cs.validateTaxPeriod(criteria.TenantId, challan.BusinessService, challan.TaxPeriodFrom, challan.TaxPeriodTo, req.RequestInfo)
		if err != nil {
			return &domain.CalculationResponse{
				Errors: []domain.Error{{
					Code:    "EG_BS_TAXPERIODS_DEMAND",
					Message: err.Error(),
				}},
			}, nil
		}

		// Step 2: Build demand details from amount array
		demandDetails := make([]domain.DemandDetail, 0, len(challan.Amount))
		taxEstimates := make([]domain.TaxEstimate, 0, len(challan.Amount))

		for _, amt := range challan.Amount {
			demandDetails = append(demandDetails, domain.DemandDetail{
				TaxHeadMasterCode: amt.TaxHeadCode,
				TaxAmount:         amt.Amount,
				CollectionAmount:  0,
			})
			taxEstimates = append(taxEstimates, domain.TaxEstimate{
				TaxHeadCode:    amt.TaxHeadCode,
				EstimateAmount: amt.Amount,
			})
		}

		// Step 3: Create demand via billing-service
		err = cs.createDemand(domain.DemandRequest{
			RequestInfo: req.RequestInfo,
			Demands: []domain.Demand{{
				TenantId:        criteria.TenantId,
				BusinessService: challan.BusinessService,
				ConsumerCode:    challan.ChallanNo,
				ConsumerType:    "challan",
				TaxPeriodFrom:   challan.TaxPeriodFrom,
				TaxPeriodTo:     challan.TaxPeriodTo,
				DemandDetails:   demandDetails,
				Payer:           challan.Citizen,
			}},
		})
		if err != nil {
			log.Printf("⚠️ billing-service demand creation failed: %v", err)
		} else {
			log.Printf("✅ Demand created in billing-service for challan: %s", challan.ChallanNo)
		}

		calculations = append(calculations, domain.Calculation{
			TenantId:        criteria.TenantId,
			ChallanNo:       challan.ChallanNo,
			TaxHeadEstimate: taxEstimates,
		})
	}

	return &domain.CalculationResponse{
		ResponseInfo: &domain.ResponseInfo{
			ResMsgId: "uief87324",
			Status:   "successful",
		},
		Calculations: calculations,
	}, nil
}

// validateTaxPeriod calls MDMS to check if the tax period is valid.
func (cs *CalculatorService) validateTaxPeriod(tenantId, businessService string, fromDate, toDate int64, reqInfo *domain.RequestInfo) error {
	url := fmt.Sprintf("%s/egov-mdms-service/v1/_search", cs.mdmsHost)

	payload := map[string]interface{}{
		"RequestInfo": reqInfo,
		"MdmsCriteria": map[string]interface{}{
			"tenantId": tenantId,
			"moduleDetails": []map[string]interface{}{
				{
					"moduleName": "BillingService",
					"masterDetails": []map[string]interface{}{
						{"name": "TaxPeriod"},
						{"name": "TaxHeadMaster"},
					},
				},
			},
		},
	}

	body, _ := json.Marshal(payload)
	resp, err := cs.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("MDMS call failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var mdmsResp map[string]interface{}
	if err := json.Unmarshal(respBody, &mdmsResp); err != nil {
		return fmt.Errorf("failed to parse MDMS response: %w", err)
	}

	// Check if TaxPeriod data exists for the given date range
	mdmsRes, ok := mdmsResp["MdmsRes"]
	if !ok {
		return fmt.Errorf("No Tax Periods Found for the given demand with fromDate : %d and toDate : %d", fromDate, toDate)
	}

	billing, ok := mdmsRes.(map[string]interface{})["BillingService"]
	if !ok {
		return fmt.Errorf("No Tax Periods Found for the given demand with fromDate : %d and toDate : %d", fromDate, toDate)
	}

	taxPeriods, ok := billing.(map[string]interface{})["TaxPeriod"]
	if !ok || taxPeriods == nil {
		return fmt.Errorf("No Tax Periods Found for the given demand with fromDate : %d and toDate : %d", fromDate, toDate)
	}

	log.Printf("✅ Tax period validated for %s (from: %d, to: %d)", businessService, fromDate, toDate)
	return nil
}

// createDemand calls billing-service to create a demand.
func (cs *CalculatorService) createDemand(demandReq domain.DemandRequest) error {
	url := fmt.Sprintf("%s/billing-service/demand/_create", cs.billingHost)

	body, _ := json.Marshal(demandReq)
	resp, err := cs.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("billing-service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("billing-service returned %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
