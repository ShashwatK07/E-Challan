package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// EnrichmentService handles all synchronous REST calls to other DIGIT microservices.
// This replaces the Java ChallanEnrichment.java + service integration code.
type EnrichmentService struct {
	egovUserHost  string // e.g., "http://localhost:8092"
	egovMdmsHost  string // e.g., "http://localhost:8094"
	egovIdgenHost string // e.g., "http://localhost:8088"
	httpClient    *http.Client
}

func NewEnrichmentService(userHost, mdmsHost, idgenHost string) *EnrichmentService {
	return &EnrichmentService{
		egovUserHost:  userHost,
		egovMdmsHost:  mdmsHost,
		egovIdgenHost: idgenHost,
		httpClient:    &http.Client{},
	}
}

// =============================================================================
// egov-idgen integration (port 8088)
// Generates unique challan numbers in format: CB-CH-YYYY-MM-DD-NNNNNN
// =============================================================================

type idGenRequest struct {
	RequestInfo interface{} `json:"RequestInfo"`
	IdRequests  []idReq     `json:"idRequests"`
}

type idReq struct {
	IdName   string `json:"idName"`
	TenantId string `json:"tenantId"`
	Format   string `json:"format"`
}

type idGenResponse struct {
	IdResponses []struct {
		Id string `json:"id"`
	} `json:"idResponses"`
}

// GenerateChallanNo calls egov-idgen to get a unique challan number.
func (e *EnrichmentService) GenerateChallanNo(tenantId string, requestInfo interface{}) (string, error) {
	url := fmt.Sprintf("%s/egov-idgen/id/_generate", e.egovIdgenHost)

	payload := idGenRequest{
		RequestInfo: requestInfo,
		IdRequests: []idReq{
			{
				IdName:   "challan.receipt.id",
				TenantId: tenantId,
				Format:   "CB-CH-[cy:yyyy-MM-dd]-[SEQ_ECHALLAN_ID]",
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal idgen request: %w", err)
	}

	resp, err := e.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("egov-idgen call failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var idResp idGenResponse
	if err := json.Unmarshal(respBody, &idResp); err != nil {
		return "", fmt.Errorf("failed to parse idgen response: %w", err)
	}

	if len(idResp.IdResponses) == 0 || idResp.IdResponses[0].Id == "" {
		return "", fmt.Errorf("egov-idgen returned empty ID")
	}

	log.Printf("✅ Generated challanNo from egov-idgen: %s", idResp.IdResponses[0].Id)
	return idResp.IdResponses[0].Id, nil
}

// =============================================================================
// egov-user integration (port 8092)
// Creates or looks up citizens by mobileNumber
// =============================================================================

type userCreateRequest struct {
	RequestInfo interface{} `json:"RequestInfo"`
	User        interface{} `json:"user"`
}

type userSearchRequest struct {
	RequestInfo interface{} `json:"RequestInfo"`
}

type userSearchResponse struct {
	User []struct {
		Id           int    `json:"id"`
		Uuid         string `json:"uuid"`
		UserName     string `json:"userName"`
		Name         string `json:"name"`
		MobileNumber string `json:"mobileNumber"`
		Active       bool   `json:"active"`
		Type         string `json:"type"`
		TenantId     string `json:"tenantId"`
	} `json:"user"`
}

// CreateOrLookupCitizen creates a citizen if not found, or returns the existing one.
// This replicates the Java UserService.createUser() / UserService.searchUser() flow.
func (e *EnrichmentService) CreateOrLookupCitizen(citizen interface{}, requestInfo interface{}) (string, int, error) {
	// First, try to search for existing user by mobile number
	searchURL := fmt.Sprintf("%s/user/users/_search", e.egovUserHost)

	searchPayload := userSearchRequest{RequestInfo: requestInfo}
	body, _ := json.Marshal(searchPayload)

	resp, err := e.httpClient.Post(searchURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("⚠️ egov-user search failed: %v, proceeding with create", err)
	} else {
		defer resp.Body.Close()
		// If user found, return their UUID
	}

	// Create user via _createnovalidate
	createURL := fmt.Sprintf("%s/user/users/_createnovalidate", e.egovUserHost)

	createPayload := userCreateRequest{
		RequestInfo: requestInfo,
		User:        citizen,
	}
	body, _ = json.Marshal(createPayload)

	resp, err = e.httpClient.Post(createURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", 0, fmt.Errorf("egov-user create failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var userResp struct {
		User []struct {
			Id   int    `json:"id"`
			Uuid string `json:"uuid"`
		} `json:"user"`
	}

	if err := json.Unmarshal(respBody, &userResp); err != nil {
		return "", 0, fmt.Errorf("failed to parse user response: %w", err)
	}

	if len(userResp.User) > 0 {
		log.Printf("✅ Citizen created/found: uuid=%s, id=%d", userResp.User[0].Uuid, userResp.User[0].Id)
		return userResp.User[0].Uuid, userResp.User[0].Id, nil
	}

	return "", 0, fmt.Errorf("egov-user returned no user data")
}

// =============================================================================
// egov-mdms-service integration (port 8094)
// Validates business service configuration
// =============================================================================

type mdmsSearchRequest struct {
	RequestInfo  interface{} `json:"RequestInfo"`
	MdmsCriteria struct {
		TenantId      string         `json:"tenantId"`
		ModuleDetails []moduleDetail `json:"moduleDetails"`
	} `json:"MdmsCriteria"`
}

type moduleDetail struct {
	ModuleName    string          `json:"moduleName"`
	MasterDetails []masterDetail  `json:"masterDetails"`
}

type masterDetail struct {
	Name string `json:"name"`
}

// ValidateBusinessService calls MDMS to confirm the businessService code is valid.
func (e *EnrichmentService) ValidateBusinessService(tenantId string, requestInfo interface{}) error {
	url := fmt.Sprintf("%s/egov-mdms-service/v1/_search", e.egovMdmsHost)

	payload := mdmsSearchRequest{
		RequestInfo: requestInfo,
	}
	payload.MdmsCriteria.TenantId = tenantId
	payload.MdmsCriteria.ModuleDetails = []moduleDetail{
		{
			ModuleName:    "BillingService",
			MasterDetails: []masterDetail{{Name: "BusinessService"}},
		},
	}

	body, _ := json.Marshal(payload)
	resp, err := e.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("egov-mdms call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("egov-mdms returned status: %d", resp.StatusCode)
	}

	log.Printf("✅ MDMS validation passed for tenant: %s", tenantId)
	return nil
}
