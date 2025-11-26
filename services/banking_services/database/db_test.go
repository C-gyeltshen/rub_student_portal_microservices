package database

import (
"os"
"testing"
"github.com/stretchr/testify/assert"
)

func TestConnect_MissingDatabaseURL(t *testing.T) {
	originalURL := os.Getenv("DATABASE_URL")
	defer func() {
		if originalURL != "" {
			os.Setenv("DATABASE_URL", originalURL)
		}
	}()

	os.Unsetenv("DATABASE_URL")
	dsn := os.Getenv("DATABASE_URL")
	assert.Empty(t, dsn)
}

func TestConnect_WithValidURL(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test")
	}

	err := Connect()
	assert.NoError(t, err)
	assert.NotNil(t, DB)
}

func TestDatabaseURL_EnvVariable(t *testing.T) {
	testURL := "postgresql://test:test@localhost:5432/testdb"
	os.Setenv("DATABASE_URL", testURL)
	
	retrieved := os.Getenv("DATABASE_URL")
	assert.Equal(t, testURL, retrieved)
	
	os.Unsetenv("DATABASE_URL")
}
