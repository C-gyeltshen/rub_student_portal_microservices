package services

import (
	"database/sql"
	"finance_service/database"
	"finance_service/models"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDB is the global test database instance
var TestDB *gorm.DB

// setupTestDB initializes the test database
func setupTestDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=rub_student_portal port=5434 sslmode=disable"
	}

	var err error
	TestDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Set the global database for the services
	database.DB = TestDB

	// Disable foreign key checks for tests
	TestDB.Exec("SET session_replication_role = replica")

	return nil
}

// teardownTestDB cleans up the test database
func teardownTestDB() {
	if TestDB != nil {
		sqlDB, _ := TestDB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}
}

// TestTransactionCreation tests basic transaction creation in database
func TestTransactionCreation(t *testing.T) {
	if err := setupTestDB(); err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer teardownTestDB()

	// Clean up test data
	TestDB.Exec("DELETE FROM transactions")

	// Create test transaction
	transactionID := uuid.New()
	stipendID := uuid.New()
	studentID := uuid.New()

	testTransaction := &models.Transaction{
		ID:                 transactionID,
		StipendID:          stipendID,
		StudentID:          studentID,
		Amount:             5000.00,
		SourceAccount:      "INSTITUTION",
		DestinationAccount: "1234567890",
		DestinationBank:    "Test Bank",
		TransactionType:    models.TransactionTypeStipend,
		Status:             models.TransactionStatusPending,
		PaymentMethod:      models.PaymentMethodBankTransfer,
		InitiatedAt:        time.Now(),
	}

	if err := TestDB.Create(testTransaction).Error; err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Verify transaction was created
	var retrievedTransaction models.Transaction
	if err := TestDB.First(&retrievedTransaction, "id = ?", transactionID).Error; err != nil {
		t.Fatalf("Failed to retrieve transaction: %v", err)
	}

	if retrievedTransaction.ID != transactionID {
		t.Errorf("Transaction ID mismatch: expected %s, got %s", transactionID, retrievedTransaction.ID)
	}

	if retrievedTransaction.Status != models.TransactionStatusPending {
		t.Errorf("Status mismatch: expected %s, got %s", models.TransactionStatusPending, retrievedTransaction.Status)
	}

	if retrievedTransaction.Amount != 5000.00 {
		t.Errorf("Amount mismatch: expected 5000.00, got %.2f", retrievedTransaction.Amount)
	}

	t.Logf("✓ Transaction created and retrieved successfully")
	t.Logf("  ID: %s", retrievedTransaction.ID)
	t.Logf("  Status: %s", retrievedTransaction.Status)
	t.Logf("  Amount: %.2f", retrievedTransaction.Amount)
}

// TestTransactionStatusUpdate tests updating transaction status
func TestTransactionStatusUpdate(t *testing.T) {
	if err := setupTestDB(); err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer teardownTestDB()

	// Clean up test data
	TestDB.Exec("DELETE FROM transactions")

	// Create test transaction
	transactionID := uuid.New()
	testTransaction := &models.Transaction{
		ID:                 transactionID,
		StipendID:          uuid.New(),
		StudentID:          uuid.New(),
		Amount:             5000.00,
		DestinationAccount: "1234567890",
		DestinationBank:    "Test Bank",
		Status:             models.TransactionStatusPending,
		PaymentMethod:      models.PaymentMethodBankTransfer,
		InitiatedAt:        time.Now(),
	}

	TestDB.Create(testTransaction)

	// Update status
	now := time.Now()
	if err := TestDB.Model(&testTransaction).Updates(map[string]interface{}{
		"status":       models.TransactionStatusProcessing,
		"processed_at": now,
	}).Error; err != nil {
		t.Fatalf("Failed to update transaction: %v", err)
	}

	// Verify update
	var updatedTransaction models.Transaction
	TestDB.First(&updatedTransaction, "id = ?", transactionID)

	if updatedTransaction.Status != models.TransactionStatusProcessing {
		t.Errorf("Status mismatch: expected %s, got %s", models.TransactionStatusProcessing, updatedTransaction.Status)
	}

	if updatedTransaction.ProcessedAt == nil {
		t.Error("ProcessedAt is nil after update")
	}

	t.Logf("✓ Transaction status updated successfully")
	t.Logf("  Old Status: %s", models.TransactionStatusPending)
	t.Logf("  New Status: %s", updatedTransaction.Status)
}

// TestTransactionSuccess tests updating transaction to success
func TestTransactionSuccess(t *testing.T) {
	if err := setupTestDB(); err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer teardownTestDB()

	// Clean up test data
	TestDB.Exec("DELETE FROM transactions")

	// Create test transaction
	transactionID := uuid.New()
	testTransaction := &models.Transaction{
		ID:                 transactionID,
		StipendID:          uuid.New(),
		StudentID:          uuid.New(),
		Amount:             5000.00,
		DestinationAccount: "1234567890",
		DestinationBank:    "Test Bank",
		Status:             models.TransactionStatusPending,
		PaymentMethod:      models.PaymentMethodBankTransfer,
		InitiatedAt:        time.Now(),
	}

	TestDB.Create(testTransaction)

	// Update to success
	now := time.Now()
	referenceNumber := fmt.Sprintf("TXN-%d-%s", time.Now().UnixNano(), transactionID.String()[:8])

	if err := TestDB.Model(&testTransaction).Updates(map[string]interface{}{
		"status":           models.TransactionStatusSuccess,
		"reference_number": referenceNumber,
		"completed_at":     now,
	}).Error; err != nil {
		t.Fatalf("Failed to update transaction: %v", err)
	}

	// Verify update
	var successTransaction models.Transaction
	TestDB.First(&successTransaction, "id = ?", transactionID)

	if successTransaction.Status != models.TransactionStatusSuccess {
		t.Errorf("Status mismatch: expected %s, got %s", models.TransactionStatusSuccess, successTransaction.Status)
	}

	if !successTransaction.ReferenceNumber.Valid || successTransaction.ReferenceNumber.String == "" {
		t.Error("ReferenceNumber is empty after successful transaction")
	}

	if successTransaction.CompletedAt == nil {
		t.Error("CompletedAt is nil after successful transaction")
	}

	t.Logf("✓ Transaction marked as success successfully")
	t.Logf("  Status: %s", successTransaction.Status)
	if successTransaction.ReferenceNumber.Valid {
		t.Logf("  Reference Number: %s", successTransaction.ReferenceNumber.String)
	}
	t.Logf("  Completed At: %s", successTransaction.CompletedAt)
}

// TestTransactionFailed tests handling failed transaction
func TestTransactionFailed(t *testing.T) {
	if err := setupTestDB(); err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer teardownTestDB()

	// Clean up test data
	TestDB.Exec("DELETE FROM transactions")

	// Create test transaction
	transactionID := uuid.New()
	testTransaction := &models.Transaction{
		ID:                 transactionID,
		StipendID:          uuid.New(),
		StudentID:          uuid.New(),
		Amount:             5000.00,
		DestinationAccount: "1234567890",
		DestinationBank:    "Test Bank",
		Status:             models.TransactionStatusProcessing,
		PaymentMethod:      models.PaymentMethodBankTransfer,
		InitiatedAt:        time.Now(),
	}

	TestDB.Create(testTransaction)

	// Update to failed
	errorMsg := "Payment gateway timeout"

	if err := TestDB.Model(&testTransaction).Updates(map[string]interface{}{
		"status":          models.TransactionStatusFailed,
		"error_message":   errorMsg,
	}).Error; err != nil {
		t.Fatalf("Failed to update transaction: %v", err)
	}

	// Verify update
	var failedTransaction models.Transaction
	TestDB.First(&failedTransaction, "id = ?", transactionID)

	if failedTransaction.Status != models.TransactionStatusFailed {
		t.Errorf("Status mismatch: expected %s, got %s", models.TransactionStatusFailed, failedTransaction.Status)
	}

	if failedTransaction.ErrorMessage != errorMsg {
		t.Errorf("Error message mismatch: expected '%s', got '%s'", errorMsg, failedTransaction.ErrorMessage)
	}

	t.Logf("✓ Transaction marked as failed successfully")
	t.Logf("  Status: %s", failedTransaction.Status)
	t.Logf("  Error: %s", failedTransaction.ErrorMessage)
}

// TestTransactionQuery tests querying transactions
func TestTransactionQuery(t *testing.T) {
	if err := setupTestDB(); err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer teardownTestDB()

	// Clean up test data
	TestDB.Exec("DELETE FROM transactions")

	// Create multiple test transactions for same stipend
	stipendID := uuid.New()
	studentID := uuid.New()

	for i := 0; i < 3; i++ {
		testTransaction := &models.Transaction{
			ID:                 uuid.New(),
			StipendID:          stipendID,
			StudentID:          studentID,
			Amount:             5000.00,
			DestinationAccount: "1234567890",
			DestinationBank:    "Test Bank",
			Status:             models.TransactionStatusPending,
			PaymentMethod:      models.PaymentMethodBankTransfer,
			ReferenceNumber:    sql.NullString{String: fmt.Sprintf("TXN-TEST-%d-%d", time.Now().UnixNano(), i), Valid: true},
			InitiatedAt:        time.Now(),
		}
		TestDB.Create(testTransaction)
	}

	// Query transactions by stipend ID
	var transactions []models.Transaction
	if err := TestDB.Where("stipend_id = ?", stipendID).Find(&transactions).Error; err != nil {
		t.Fatalf("Failed to query transactions: %v", err)
	}

	if len(transactions) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(transactions))
	}

	// Query transactions by student ID
	var studentTransactions []models.Transaction
	if err := TestDB.Where("student_id = ?", studentID).Find(&studentTransactions).Error; err != nil {
		t.Fatalf("Failed to query transactions: %v", err)
	}

	if len(studentTransactions) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(studentTransactions))
	}

	t.Logf("✓ Transaction queries successful")
	t.Logf("  Found %d transactions for stipend", len(transactions))
	t.Logf("  Found %d transactions for student", len(studentTransactions))
}

// TestTransactionCancellation tests canceling a transaction
func TestTransactionCancellation(t *testing.T) {
	if err := setupTestDB(); err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer teardownTestDB()

	// Clean up test data
	TestDB.Exec("DELETE FROM transactions")

	// Create test transaction
	transactionID := uuid.New()
	testTransaction := &models.Transaction{
		ID:                 transactionID,
		StipendID:          uuid.New(),
		StudentID:          uuid.New(),
		Amount:             5000.00,
		DestinationAccount: "1234567890",
		DestinationBank:    "Test Bank",
		Status:             models.TransactionStatusPending,
		PaymentMethod:      models.PaymentMethodBankTransfer,
		InitiatedAt:        time.Now(),
	}

	TestDB.Create(testTransaction)

	// Cancel transaction
	cancelReason := "User requested cancellation"

	if err := TestDB.Model(&testTransaction).Updates(map[string]interface{}{
		"status":   models.TransactionStatusCancelled,
		"remarks":  cancelReason,
	}).Error; err != nil {
		t.Fatalf("Failed to cancel transaction: %v", err)
	}

	// Verify cancellation
	var cancelledTransaction models.Transaction
	TestDB.First(&cancelledTransaction, "id = ?", transactionID)

	if cancelledTransaction.Status != models.TransactionStatusCancelled {
		t.Errorf("Status mismatch: expected %s, got %s", models.TransactionStatusCancelled, cancelledTransaction.Status)
	}

	if cancelledTransaction.Remarks != cancelReason {
		t.Errorf("Remarks mismatch: expected '%s', got '%s'", cancelReason, cancelledTransaction.Remarks)
	}

	t.Logf("✓ Transaction cancelled successfully")
	t.Logf("  Status: %s", cancelledTransaction.Status)
	t.Logf("  Reason: %s", cancelledTransaction.Remarks)
}
