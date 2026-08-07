package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/CDPI-HRSS/echallan-calculator/internal/domain"
	"github.com/CDPI-HRSS/echallan-calculator/internal/service"
)

// CalculatorController handles the HTTP endpoint for the eChallan Calculator.
// This is the Go equivalent of Java ChallanCalController.java.
type CalculatorController struct {
	calcService *service.CalculatorService
}

func NewCalculatorController(svc *service.CalculatorService) *CalculatorController {
	return &CalculatorController{calcService: svc}
}

// RegisterRoutes sets up the calculator route.
// Path: /echallan-calculator/v1/_calculate (matches Java @RequestMapping("/v1"))
func (ctrl *CalculatorController) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/echallan-calculator/v1")
	{
		v1.POST("/_calculate", ctrl.Calculate)
	}
}

// Calculate handles POST /echallan-calculator/v1/_calculate
// Java method: ChallanCalController.calculate()
// E2E report Step 9.
func (ctrl *CalculatorController) Calculate(c *gin.Context) {
	var req domain.CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "INVALID_REQUEST",
				Message: "Failed to parse JSON payload: " + err.Error(),
			}},
		})
		return
	}

	result, err := ctrl.calcService.Calculate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: []domain.Error{{
				Code:    "CALCULATION_ERROR",
				Message: err.Error(),
			}},
		})
		return
	}

	// If the result contains errors (e.g., missing TaxPeriod), return them
	if len(result.Errors) > 0 {
		c.JSON(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusOK, result)
}
