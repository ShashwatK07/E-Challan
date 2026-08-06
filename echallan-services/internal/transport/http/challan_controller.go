package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/CDPI-HRSS/echallan-services/internal/domain"
	"github.com/CDPI-HRSS/echallan-services/internal/service"
)

// ChallanController handles all HTTP endpoints for the eChallan service.
// This is the Go equivalent of Java ChallanController.java.
type ChallanController struct {
	challanService service.ChallanService
}

func NewChallanController(svc service.ChallanService) *ChallanController {
	return &ChallanController{challanService: svc}
}

// RegisterRoutes sets up all routes on the Gin engine.
// The context path /echallan-services is set at the engine level in main.go.
// These routes replicate the Java @RequestMapping("eChallan/v1") structure.
func (cc *ChallanController) RegisterRoutes(router *gin.Engine) {
	// Legacy DIGIT routes — exact path parity with Java (E2E report)
	v1 := router.Group("/echallan-services/eChallan/v1")
	{
		v1.POST("/_create", cc.Create)   // Step 4
		v1.POST("/_search", cc.Search)   // Step 5
		v1.POST("/_update", cc.Update)   // Step 6
		v1.POST("/_count", cc.Count)     // Step 7
		v1.POST("/_test", cc.Test)       // Step 8
	}
}

// Create handles POST /echallan-services/eChallan/v1/_create
// Java method: ChallanController.create()
func (cc *ChallanController) Create(c *gin.Context) {
	var req domain.ChallanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "INVALID_REQUEST",
				Message: "Failed to parse JSON payload: " + err.Error(),
			}},
		})
		return
	}

	challan, err := cc.challanService.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "CREATE_ERROR",
				Message: err.Error(),
			}},
		})
		return
	}

	// Build response matching EXACT Java output (E2E report Step 4)
	res := domain.ChallanResponse{
		CountOfServices:      0,
		TotalAmountCollected: 0,
		ChallanValidity:      0,
		ResponseInfo:         createResponseInfo(req.RequestInfo, "successful"),
		Challans:             []*domain.Challan{challan},
		TotalCount:           0,
	}
	c.JSON(http.StatusOK, res)
}

// Search handles POST /echallan-services/eChallan/v1/_search
// Java method: ChallanController.search()
func (cc *ChallanController) Search(c *gin.Context) {
	var reqWrapper domain.RequestInfoWrapper
	if err := c.ShouldBindJSON(&reqWrapper); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "INVALID_REQUEST",
				Message: "Failed to parse JSON payload: " + err.Error(),
			}},
		})
		return
	}

	// Search criteria come from query parameters (E2E report Step 5)
	var criteria domain.SearchCriteria
	if err := c.ShouldBindQuery(&criteria); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "INVALID_QUERY",
				Message: "Failed to parse query parameters: " + err.Error(),
			}},
		})
		return
	}

	challans, err := cc.challanService.Search(criteria, reqWrapper.RequestInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "SEARCH_ERROR",
				Message: err.Error(),
			}},
		})
		return
	}

	// Count unique business services for the response metadata
	serviceSet := make(map[string]bool)
	validCount := 0
	for _, ch := range challans {
		serviceSet[ch.BusinessService] = true
		if ch.ApplicationStatus == "ACTIVE" {
			validCount++
		}
	}

	res := domain.ChallanResponse{
		CountOfServices:      len(serviceSet),
		TotalAmountCollected: 0,
		ChallanValidity:      validCount,
		ResponseInfo:         createResponseInfo(reqWrapper.RequestInfo, "successful"),
		Challans:             challans,
		TotalCount:           len(challans),
	}
	c.JSON(http.StatusOK, res)
}

// Update handles POST /echallan-services/eChallan/v1/_update
// Java method: ChallanController.update()
func (cc *ChallanController) Update(c *gin.Context) {
	var req domain.ChallanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "INVALID_REQUEST",
				Message: "Failed to parse JSON payload: " + err.Error(),
			}},
		})
		return
	}

	challan, err := cc.challanService.Update(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "UPDATE_ERROR",
				Message: err.Error(),
			}},
		})
		return
	}

	res := domain.ChallanResponse{
		CountOfServices:      0,
		TotalAmountCollected: 0,
		ChallanValidity:      0,
		ResponseInfo:         createResponseInfo(req.RequestInfo, "successful"),
		Challans:             []*domain.Challan{challan},
		TotalCount:           0,
	}
	c.JSON(http.StatusOK, res)
}

// Count handles POST /echallan-services/eChallan/v1/_count
// Java method: ChallanController.count()
// NOTE: Returns a DIFFERENT response structure than Create/Search/Update!
// Uses uppercase "ResponseInfo" and "ChallanCount" (E2E report Step 7).
func (cc *ChallanController) Count(c *gin.Context) {
	tenantId := c.Query("tenantId")

	var reqWrapper domain.RequestInfoWrapper
	if err := c.ShouldBindJSON(&reqWrapper); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "INVALID_REQUEST",
				Message: "Failed to parse JSON payload: " + err.Error(),
			}},
		})
		return
	}

	countResp, err := cc.challanService.Count(tenantId, reqWrapper.RequestInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "COUNT_ERROR",
				Message: err.Error(),
			}},
		})
		return
	}

	c.JSON(http.StatusOK, countResp)
}

// Test handles POST /echallan-services/eChallan/v1/_test
// Java method: ChallanController.test()
// Simply pushes to Kafka to verify the producer is working (E2E report Step 8).
func (cc *ChallanController) Test(c *gin.Context) {
	var req domain.ChallanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "INVALID_REQUEST",
				Message: "Failed to parse JSON payload: " + err.Error(),
			}},
		})
		return
	}

	if err := cc.challanService.Test(&req); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "KAFKA_TEST_ERROR",
				Message: err.Error(),
			}},
		})
		return
	}

	res := domain.ChallanResponse{
		ResponseInfo: createResponseInfo(req.RequestInfo, "successful"),
		Challans:     []*domain.Challan{req.Challan},
	}
	c.JSON(http.StatusOK, res)
}

// createResponseInfo builds the standard DIGIT ResponseInfo from the incoming RequestInfo.
func createResponseInfo(reqInfo *domain.RequestInfo, status string) *domain.ResponseInfo {
	if reqInfo == nil {
		return &domain.ResponseInfo{
			ResMsgId: "uief87324",
			Status:   status,
		}
	}
	return &domain.ResponseInfo{
		ApiId:    reqInfo.ApiId,
		Ver:      reqInfo.Ver,
		ResMsgId: "uief87324",
		MsgId:    reqInfo.MsgId,
		Status:   status,
	}
}
