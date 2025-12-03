package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect_MissingDatabaseURL(t *testing.T) {
	os.Unsetenv("DATABASE_URL")
	
	err := Connect()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid db")
}

func TestDatabaseURL_EnvVariable(t *testing.T) {
	testURL := "postgres://testuser:testpass@localhost/testdb"
	os.Setenv("DATABASE_URL", testURL)
	defer os.Unsetenv("DATABASE_URL")
	
	url := os.Getenv("DATABASE_URL")
	assert.Equal(t, testURL, url)
}
