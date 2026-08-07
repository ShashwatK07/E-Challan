package configs

import "os"

// Config holds all configuration for the echallan-calculator application.
type Config struct {
	ServerPort  string
	ContextPath string

	// Dependent Service URLs
	EgovMdmsHost    string // egov-mdms-service (port 8094)
	BillingHost     string // billing-service (port 8081)
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8078"),
		ContextPath: getEnv("CONTEXT_PATH", "/echallan-calculator"),

		EgovMdmsHost:    getEnv("EGOV_MDMS_HOST", "http://localhost:8094"),
		BillingHost:     getEnv("BILLING_HOST", "http://localhost:8081"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
