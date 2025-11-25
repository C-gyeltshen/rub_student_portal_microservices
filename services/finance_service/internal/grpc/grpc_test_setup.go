package grpc

import (
	"fmt"
	"os"
	"testing"
	"time"

	"finance_service/database"

	"github.com/google/uuid"
)

var testStudentID string

// GetTestStudentID returns a valid test student ID, initializing it if needed
func GetTestStudentID() string {
	if testStudentID != "" {
		return testStudentID
	}

	// Ensure database is connected
	if database.DB == nil {
		if err := database.Connect(); err != nil {
			fmt.Printf("Error connecting to database: %v\n", err)
			return uuid.New().String()
		}
	}

	// Initialize it now if not already done
	if err := initializeTestStudent(); err != nil {
		fmt.Printf("Warning: Failed to initialize test student: %v\n", err)
		// Return a fallback UUID
		return uuid.New().String()
	}

	return testStudentID
}

// TestMain initializes test infrastructure for gRPC tests
func TestMain(m *testing.M) {
	fmt.Println("========== TestMain called ==========")
	// Set DATABASE_URL if not already set
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "postgresql://postgres:postgres@postgres:5432/rub_student_portal?sslmode=disable")
	}

	// Initialize database connection
	if err := database.Connect(); err != nil {
		panic("Failed to initialize test database: " + err.Error())
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		panic("Failed to run migrations: " + err.Error())
	}

	// Initialize database
	if err := database.InitializeFinanceDatabase(); err != nil {
		fmt.Printf("Warning: Failed to initialize finance database: %v\n", err)
	}

	// Clean up old test data before running tests
	if err := cleanupOldTestData(); err != nil {
		fmt.Printf("Warning: Failed to cleanup old test data: %v\n", err)
	}

	// Initialize test student ID by fetching or creating a student
	if err := initializeTestStudent(); err != nil {
		fmt.Printf("Warning: Failed to initialize test student: %v\n", err)
		// Fallback to a default UUID
		testStudentID = uuid.New().String()
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

// cleanupOldTestData removes old test deduction rules to prevent duplicates
func cleanupOldTestData() error {
	if database.DB == nil {
		return nil
	}

	// Delete test rules older than 1 hour
	cutoffTime := time.Now().Add(-1 * time.Hour)
	return database.DB.Where("rule_name LIKE ? AND created_at < ?", "Test%", cutoffTime).
		Delete(nil).Error
}

// GetUniqueTestName generates a unique test name with timestamp to avoid duplicates
func GetUniqueTestName(baseName string) string {
	return fmt.Sprintf("%s_%d", baseName, time.Now().UnixNano())
}

// initializeTestStudent retrieves an existing test student or creates one for testing
func initializeTestStudent() error {
	if database.DB == nil {
		fmt.Println("Warning: database.DB is nil, using fallback student ID")
		testStudentID = uuid.New().String()
		return nil
	}

	// Try to get an existing test student from the students table
	var studentID string
	result := database.DB.Raw("SELECT id FROM students LIMIT 1").Scan(&studentID)
	
	if result.Error == nil && studentID != "" {
		testStudentID = studentID
		fmt.Printf("Using existing test student ID: %s\n", testStudentID)
		return nil
	}

	// No existing student found, create one for testing
	testID := uuid.New().String()
	// Retry until we find a unique email and card number
	for attempts := 0; attempts < 5; attempts++ {
		email := fmt.Sprintf("test-student-%d@rub.edu.bt", time.Now().UnixNano())
		cardNum := fmt.Sprintf("TEST%d", time.Now().UnixNano())
		
		// Try to insert directly without conflict clause
		if err := database.DB.Exec(
			`INSERT INTO students (id, email, name, rub_id_card_number, phone_number)
			 VALUES ($1, $2, $3, $4, $5)`,
			testID,
			email,
			"Test Student",
			cardNum,
			"+97517123456",
		).Error; err == nil {
			// Success!
			testStudentID = testID
			fmt.Printf("Created test student ID: %s\n", testStudentID)
			return nil
		} else {
			// Try again with a different ID
			testID = uuid.New().String()
		}
		time.Sleep(10 * time.Millisecond)
	}

	return fmt.Errorf("failed to create test student after retries")
}

// TestInitialization ensures TestMain was called and student ID is initialized
func TestInitialization(t *testing.T) {
	if testStudentID == "" {
		t.Fatalf("testStudentID was not initialized. TestMain may not have been called")
	}
	t.Logf("Test student ID is initialized: %s", testStudentID)
}
