package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/CDPI-HRSS/echallan-calculator/configs"
	"github.com/CDPI-HRSS/echallan-calculator/internal/service"
	httpTransport "github.com/CDPI-HRSS/echallan-calculator/internal/transport/http"
)

// main is the entry point for the Go eChallan Calculator application.
// This replaces the Java Spring Boot ChallanCalculationApplication.java.
func main() {
	log.Println("🚀 Starting eChallan Calculator (Go) ...")

	// 1. Load configuration
	cfg := configs.LoadConfig()
	log.Printf("📋 Config loaded: port=%s, mdms=%s, billing=%s",
		cfg.ServerPort, cfg.EgovMdmsHost, cfg.BillingHost)

	// 2. Initialize calculator service
	calcService := service.NewCalculatorService(cfg.EgovMdmsHost, cfg.BillingHost)

	// 3. Initialize HTTP controller
	controller := httpTransport.NewCalculatorController(calcService)

	// 4. Create Gin router
	router := gin.Default()
	controller.RegisterRoutes(router)

	// 5. Start HTTP server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🌐 eChallan Calculator listening on %s", addr)
	log.Printf("📡 API endpoint: http://localhost:%s/echallan-calculator/v1/_calculate", cfg.ServerPort)
	if err := router.Run(addr); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
