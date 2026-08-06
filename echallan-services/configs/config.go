package configs

import "os"

// Config holds all configuration for the echallan-services application.
// These values map 1:1 to the Java application.properties file.
type Config struct {
	// Server
	ServerPort    string
	ContextPath   string

	// Database (PostgreSQL)
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Kafka
	KafkaBootstrapServer string
	SaveChallanTopic     string
	UpdateChallanTopic   string

	// Dependent Service URLs (synchronous REST calls)
	EgovUserHost   string // egov-user (port 8092)
	EgovMdmsHost   string // egov-mdms-service (port 8094)
	EgovIdgenHost  string // egov-idgen (port 8088)
}

// LoadConfig reads configuration from environment variables with sensible defaults
// that match the Java application.properties from the Azure VM (spiderverse).
func LoadConfig() *Config {
	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8079"),
		ContextPath:   getEnv("CONTEXT_PATH", "/echallan-services"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "devdb"),

		KafkaBootstrapServer: getEnv("KAFKA_BOOTSTRAP_SERVER", "localhost:9092"),
		SaveChallanTopic:     getEnv("SAVE_CHALLAN_TOPIC", "save-challan"),
		UpdateChallanTopic:   getEnv("UPDATE_CHALLAN_TOPIC", "update-challan"),

		EgovUserHost:   getEnv("EGOV_USER_HOST", "http://localhost:8092"),
		EgovMdmsHost:   getEnv("EGOV_MDMS_HOST", "http://localhost:8094"),
		EgovIdgenHost:  getEnv("EGOV_IDGEN_HOST", "http://localhost:8088"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
