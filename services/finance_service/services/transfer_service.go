package services

import (
	"database/sql"
	"finance_service/database"
	"finance_service/models"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// TransferService handles money transfer operations
type TransferService struct {
	bankingClient *BankingServiceClient
}

// NewTransferService creates a new transfer service
func NewTransferService() *TransferService {
	return &TransferService{
		bankingClient: NewBankingServiceClient(),
	}
}

// InitiateTransfer initiates a money transfer for a stipend
func (ts *TransferService) InitiateTransfer(stipendID uuid.UUID, paymentMethod string) (*models.Transaction, error) {
	log.Printf("Initiating transfer for stipend: %s", stipendID)

	// Fetch stipend details
	var stipend models.Stipend
	if err := database.DB.First(&stipend, "id = ?", stipendID).Error; err != nil {
		return nil, fmt.Errorf("stipend not found: %w", err)
	}

	// Check if stipend amount is greater than 0
	if stipend.Amount <= 0 {
		return nil, fmt.Errorf("invalid stipend amount: %f", stipend.Amount)
	}

	// Fetch student bank details from database (primary) or banking service (fallback)
	var bankDetails struct {
		AccountNumber string
		BankID        string
	}
	
	// Try to get from database first
	dbErr := database.DB.Table("student_bank_details").
		Select("account_number, COALESCE(bank_id::text, 'DEFAULT_BANK') as bank_id").
		Where("student_id = ?", stipend.StudentID).
		First(&bankDetails).Error
	
	// If not found in database, try banking service
	if dbErr != nil {
		log.Printf("Bank details not in database, attempting to fetch from banking service: %v", dbErr)
		bankDetailsFromService, err := ts.bankingClient.GetStudentBankDetails(stipend.StudentID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to fetch student bank details: %w", err)
		}
		bankDetails.AccountNumber = bankDetailsFromService.AccountNumber
		bankDetails.BankID = bankDetailsFromService.BankID
	}

	// Create transaction record
	transaction := &models.Transaction{
		ID:                 uuid.New(),
		StipendID:          stipendID,
		StudentID:          stipend.StudentID,
		Amount:             stipend.Amount,
		SourceAccount:      "INSTITUTION_ACCOUNT", // Placeholder - should be from config
		DestinationAccount: bankDetails.AccountNumber,
		DestinationBank:    bankDetails.BankID,
		TransactionType:    models.TransactionTypeStipend,
		Status:             models.TransactionStatusPending,
		PaymentMethod:      paymentMethod,
		ReferenceNumber:    sql.NullString{}, // Will be set during processing (NULL initially)
		InitiatedAt:        time.Now(),
	}

	// Save transaction record
	if err := database.DB.Create(transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	log.Printf("Transaction created: %s with status: %s", transaction.ID, transaction.Status)

	return transaction, nil
}

// ProcessTransfer processes a pending transfer (simulated payment gateway call)
func (ts *TransferService) ProcessTransfer(transactionID uuid.UUID) (*models.Transaction, error) {
	log.Printf("Processing transfer: %s", transactionID)

	var transaction models.Transaction
	if err := database.DB.First(&transaction, "id = ?", transactionID).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	// Check if transaction is already processed
	if transaction.Status != models.TransactionStatusPending {
		return nil, fmt.Errorf("transaction status is %s, can only process PENDING transactions", transaction.Status)
	}

	// Update status to PROCESSING
	transaction.Status = models.TransactionStatusProcessing
	now := time.Now()
	transaction.ProcessedAt = &now

	if err := database.DB.Save(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Simulate payment gateway processing
	success, referenceNumber, errorMsg := ts.simulatePaymentGatewayCall(&transaction)

	if success {
		transaction.Status = models.TransactionStatusSuccess
		transaction.ReferenceNumber = sql.NullString{String: referenceNumber, Valid: true}
		completedAt := time.Now()
		transaction.CompletedAt = &completedAt

		// Update stipend payment status
		if err := database.DB.Model(&models.Stipend{}).
			Where("id = ?", transaction.StipendID).
			Updates(map[string]interface{}{
				"payment_status": "Processed",
				"payment_date":   now,
			}).Error; err != nil {
			log.Printf("Warning: Failed to update stipend payment status: %v", err)
		}

		log.Printf("Transfer succeeded: %s with reference: %s", transactionID, referenceNumber)
	} else {
		transaction.Status = models.TransactionStatusFailed
		transaction.ErrorMessage = errorMsg
		log.Printf("Transfer failed: %s with error: %s", transactionID, errorMsg)
	}

	if err := database.DB.Save(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return &transaction, nil
}

// GetTransactionStatus retrieves the status of a transaction
func (ts *TransferService) GetTransactionStatus(transactionID uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := database.DB.First(&transaction, "id = ?", transactionID).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	return &transaction, nil
}

// GetTransactionsByStipend retrieves all transactions for a stipend
func (ts *TransferService) GetTransactionsByStipend(stipendID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := database.DB.Where("stipend_id = ?", stipendID).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	return transactions, nil
}

// GetTransactionsByStudent retrieves all transactions for a student
func (ts *TransferService) GetTransactionsByStudent(studentID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := database.DB.Where("student_id = ?", studentID).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	return transactions, nil
}

// CancelTransfer cancels a pending transaction
func (ts *TransferService) CancelTransfer(transactionID uuid.UUID, reason string) (*models.Transaction, error) {
	log.Printf("Cancelling transfer: %s", transactionID)

	var transaction models.Transaction
	if err := database.DB.First(&transaction, "id = ?", transactionID).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	// Only pending or processing transactions can be cancelled
	if transaction.Status != models.TransactionStatusPending && transaction.Status != models.TransactionStatusProcessing {
		return nil, fmt.Errorf("can only cancel PENDING or PROCESSING transactions, current status: %s", transaction.Status)
	}

	transaction.Status = models.TransactionStatusCancelled
	transaction.Remarks = reason

	if err := database.DB.Save(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to cancel transaction: %w", err)
	}

	log.Printf("Transfer cancelled: %s", transactionID)
	return &transaction, nil
}

// simulatePaymentGatewayCall simulates a call to a payment gateway
// In production, this would call actual payment gateway APIs (Stripe, PayPal, etc.)
func (ts *TransferService) simulatePaymentGatewayCall(transaction *models.Transaction) (success bool, referenceNumber string, errorMsg string) {
	// Simulate 95% success rate for demo purposes
	// In production, call actual payment gateway
	log.Printf("Simulating payment gateway call for transaction: %s, amount: %.2f", transaction.ID, transaction.Amount)

	// TODO: Replace with actual payment gateway integration
	// Example:
	// - Stripe: stripe.com
	// - PayPal: paypal.com
	// - Local bank API integration

	// For now, simulate success with 95% probability
	if transaction.Amount > 0 && transaction.Amount <= 1000000 { // Arbitrary limit for demo
		referenceNumber = fmt.Sprintf("TXN-%d-%s", time.Now().UnixNano(), transaction.ID.String()[:8])
		return true, referenceNumber, ""
	}

	return false, "", "Transfer amount exceeds limit or invalid"
}

// RetryFailedTransfer retries a failed transfer
func (ts *TransferService) RetryFailedTransfer(transactionID uuid.UUID) (*models.Transaction, error) {
	log.Printf("Retrying failed transfer: %s", transactionID)

	var transaction models.Transaction
	if err := database.DB.First(&transaction, "id = ?", transactionID).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	if transaction.Status != models.TransactionStatusFailed {
		return nil, fmt.Errorf("can only retry FAILED transactions, current status: %s", transaction.Status)
	}

	// Reset to pending before reprocessing
	transaction.Status = models.TransactionStatusPending
	transaction.ErrorMessage = ""
	if err := database.DB.Save(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to reset transaction status: %w", err)
	}

	// Try again
	return ts.ProcessTransfer(transactionID)
}
