package services

import (
	"os"
	"testing"

	"finance_service/database"
)

// TestMain is called before tests run and initializes the test database
func TestMain(m *testing.M) {
	// Set DATABASE_URL if not already set
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/rub_student_portal?sslmode=disable")
	}

	// Initialize database connection
	if err := database.Connect(); err != nil {
		panic("Failed to initialize test database: " + err.Error())
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		panic("Failed to run migrations: " + err.Error())
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}
