package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	APIEndpoint      string
	ClientID         string
	ClientSecret     string
	MongoURI         string
	MongoDatabase    string
	CollectionPrefix string
	JobSchedule      string
	Languages        []string
	DefaultLanguage  string
	LogLevel         string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (Config, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using environment variables")
	}

	// Setup configuration
	config := Config{
		APIEndpoint:      getEnv("API_ENDPOINT", "https://auth.opendatahub.testingmachine.eu/auth/realms/noi/protocol/openid-connect/token"),
		ClientID:         getEnv("CLIENT_ID", "odh-mobility-datacollector-bike-boxes"),
		ClientSecret:     getEnv("CLIENT_SECRET", "7bd46f8f-c296-416d-a13d-dc81e68d0830"),
		MongoURI:         getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:    getEnv("MONGO_DATABASE", "bikeboxes"),
		CollectionPrefix: getEnv("COLLECTION_PREFIX", "raw"),
		JobSchedule:      getEnv("JOB_SCHEDULE", "*/15 * * * *"),
		DefaultLanguage:  getEnv("DEFAULT_LANGUAGE", "it"),
		Languages:        getLanguages(),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}

	// Validate required fields
	if config.ClientID == "" || config.ClientSecret == "" {
		return config, ErrMissingCredentials
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getLanguages gets the languages from environment variables
func getLanguages() []string {
	languagesStr := getEnv("LANGUAGES", "it,en,de,lld")
	if languagesStr == "" {
		return []string{"it", "en", "de", "lld"}
	}
	
	// Split comma-separated string and trim spaces
	languages := strings.Split(languagesStr, ",")
	for i := range languages {
		languages[i] = strings.TrimSpace(languages[i])
	}
	
	return languages
}