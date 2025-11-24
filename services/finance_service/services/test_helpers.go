package services

import (
"testing"

"github.com/google/uuid"
"finance_service/database"
"finance_service/models"
)

// setupTestDBNew initializes database for tests if not already done
func setupTestDBNew(t *testing.T) {
	if database.DB == nil {
		// Try to initialize database if not already done
		if err := database.Connect(); err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}
		if err := database.Migrate(); err != nil {
			t.Fatalf("Failed to run migrations: %v", err)
		}
	}
}

// getTestStudentID returns a valid test student ID from the database
// It uses the stipend table to find a valid student
func getTestStudentID() uuid.UUID {
	// Try to get an existing student from stipends table
	var stipend models.Stipend
	if err := database.DB.First(&stipend).Error; err == nil {
		// StudentID is uuid.UUID
		if stipend.StudentID != uuid.Nil {
			return stipend.StudentID
		}
	}

	// If no stipend exists, return a fixed UUID for testing
	return uuid.MustParse("12345678-1234-1234-1234-123456789012")
}
