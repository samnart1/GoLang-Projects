package config

import (
	"os"
	"testing"

	"github.com/golang/011api/internal/config"
)

func TestLoad_WithDefaults(t *testing.T) {
	clearEnv()

	os.Setenv("DB_PASSWORD", "testpassword")
	os.Setenv("JWT_SECRET", "a")
	defer clearEnv()

	config, err := config.Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Server.Port != "8080" {
		t.Errorf("Expected port 8080, go %s", config.Server.Port)
	}

	if config.Database.Host != "localhost" {
		t.Errorf("Expected database host localhost, got %s", config.Database.Host)
	}

	if config.JWT.Secret != "a" {
		t.Errorf("Expected JWT secret a, got %s", config.JWT.Secret)
	}
}

func TestLoad_WithCustomValues(t *testing.T) {
	clearEnv()

	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_HOST", "customhost")
	os.Setenv("DB_PASSWORD", "custompassword")
	os.Setenv("JWT_SECRET", "custom-secret")
	os.Setenv("REDIS_DB", "2")
	defer clearEnv()

	config, err := config.Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Server.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", config.Server.Port)
	}

	if config.Database.Host != "customhost" {
		t.Errorf("Expected db host customhot, got %s", config.Database.Host)
	}

	if config.Redis.DB != 2 {
		t.Errorf("Expeced redis db 2, got %d", config.Redis.DB)
	}
}

func TestLoad_MissingRequiredFields(t *testing.T) {
	clearEnv()

	os.Setenv("JWT_SECRET", "test-secret")
	_, err := config.Load()
	if err == nil {
		t.Error("Expected error for missing DB_PASSWORD")
	}

	clearEnv()
	os.Setenv("DB_PASSWORD", "password")
	_, err = config.Load()
	if err == nil {
		t.Error("Expected error for missing JWT_SECRET")
	}

	clearEnv()
}

func TestGetDatabaseURL(t *testing.T) {
	config := &config.Config{
		Database: config.DatabaseConfig{
			Host: 		"localhost",
			Port: 		"5432",
			Username: 	"testuser",
			Password: 	"testpass",
			Database: 	"testdb",
			SSLMode: 	"disable",
		},
	}

	expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	actual := config.GetDatabaseURL()

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestGetServerAddress(t *testing.T) {
	config := &config.Config{
		Server: config.ServerConfig{
			Host: "0.0.0.0",
			Port: "8080",
		},
	}

	expected := "0.0.0.0:8080"
	actual := config.GetServerAddress()

	if expected != actual {
		t.Errorf("Expected %s, got %s", actual, expected)
	}
}

func TestGetRedisAddress(t *testing.T) {
	config := &config.Config{
		Redis: config.RedisConfig{
			Host: "localhost",
			Port: "6379",
		},
	}

	expected := "localhost:6379"
	actual := config.GetRedisAddress()

	if expected != actual {
		t.Errorf("Expected %s, got %s", actual, expected)
	}

}

// func TestGetDurationEnv(t *testing.T) {
// }

// func TestGetIntEnv(t *testing.T) {
	
// }

func clearEnv() {
	envVars := []string{
		"SERVER_PORT", "SERVER_HOST", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT", "DB_HOST", "DB_PORT", "DB_PASSWORD", "DB_USERNAME", "DB_SSL_MODE", "DB_NAME", "REDIS_PORT", "REDIS_HOST", "REDIS_PASSWORD", "REDIS_DB", "JWT_SECRET", "JWT_EXPIRATION",
	}

	for _, env := range envVars {
		os.Unsetenv(env)
	}
}