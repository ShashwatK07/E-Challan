package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/CDPI-HRSS/echallan-services/configs"
	"github.com/CDPI-HRSS/echallan-services/internal/repository/postgres"
	"github.com/CDPI-HRSS/echallan-services/internal/service"
	httpTransport "github.com/CDPI-HRSS/echallan-services/internal/transport/http"
	"github.com/CDPI-HRSS/echallan-services/internal/transport/kafka"
	"github.com/CDPI-HRSS/echallan-services/internal/validator"
)

// main is the entry point for the Go eChallan Services application.
// This replaces the Java Spring Boot EchallanApplication.java.
//
// Architecture flow:
//   main.go → config → repository → validator → enrichment → service → controller → gin.Engine
func main() {
	log.Println("🚀 Starting eChallan Services (Go) ...")

	// 1. Load configuration (replaces application.properties)
	cfg := configs.LoadConfig()
	log.Printf("📋 Config loaded: port=%s, db=%s:%s/%s, kafka=%s",
		cfg.ServerPort, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.KafkaBootstrapServer)

	// 2. Initialize PostgreSQL repository (replaces Spring DataSource + JdbcTemplate)
	repo, err := postgres.NewChallanRepository(
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}

	// 3. Initialize Kafka producer (replaces Spring KafkaTemplate)
	producer := kafka.NewProducer(cfg.KafkaBootstrapServer)
	defer producer.Close()

	// 4. Initialize validator
	val := validator.NewChallanValidator(repo)

	// 5. Initialize enrichment service (REST clients to egov-user, egov-mdms, egov-idgen)
	enrichment := service.NewEnrichmentService(
		cfg.EgovUserHost,
		cfg.EgovMdmsHost,
		cfg.EgovIdgenHost,
	)

	// 6. Initialize business logic service
	challanService := service.NewChallanService(repo, producer, val, enrichment)

	// 7. Initialize HTTP controller and register routes
	controller := httpTransport.NewChallanController(challanService)

	// 8. Create Gin router (replaces Spring Boot embedded Tomcat)
	router := gin.Default()
	controller.RegisterRoutes(router)

	// 9. Start HTTP server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("🌐 eChallan Services listening on %s", addr)
	log.Printf("📡 API endpoints available at http://localhost:%s/echallan-services/eChallan/v1/", cfg.ServerPort)
	if err := router.Run(addr); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
