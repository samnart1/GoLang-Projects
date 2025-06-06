package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

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

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using environment variables")
	}

	config := Config{
		APIEndpoint:      getEnv("API_ENDPOINT", ""),
		ClientID:         getEnv("CLIENT_ID", ""),
		ClientSecret:     getEnv("CLIENT_SECRET", ""),
		MongoURI:         getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:    getEnv("MONGO_DATABASE", "bikeboxes"),
		CollectionPrefix: getEnv("COLLECTION_PREFIX", "raw"),
		JobSchedule:      getEnv("JOB_SCHEDULE", "*/15 * * * *"),
		DefaultLanguage:  getEnv("DEFAULT_LANGUAGE", "it"),
		Languages:        getLanguages(),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}

	if config.ClientID == "" || config.ClientSecret == "" {
		return config, ErrMissingCredentials
	}

	return config, nil
}

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
	
	languages := strings.Split(languagesStr, ",")
	for i := range languages {
		languages[i] = strings.TrimSpace(languages[i])
	}
	
	return languages
}