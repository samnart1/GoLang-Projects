package integration

import (
	"os"
	"strconv"
	"testing"

	"github.com/golang/011api/internal/infrastructure/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration testing")
	}

	cfg := database.Config{
		Driver: 	"pgx",
		Host: 		getEnv("DB_HOST", "localhost"),
		Port: 		getIntEnv("DB_PORT", 5432),
		DBName: 	getEnv("DB_NAME", "blog_db"),
		Username: 	getEnv("DB_USER", "samnart"),
		Password: 	getEnv("DB_PASS", "Aa1234!@"),
		SSLMode: 	getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewConnection(cfg)
	require.NoError(t, err)
	defer db.Close()

	err = db.Health()
	assert.NoError(t, err)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value != "" {
		IntValue, err := strconv.Atoi(value)
		if err != nil {
			return IntValue
		}
	}
	return defaultValue
}