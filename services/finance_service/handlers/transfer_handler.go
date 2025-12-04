package handlers

import (
	"encoding/json"
	"finance_service/models"
	"finance_service/services"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TransferHandler handles HTTP requests for money transfers
type TransferHandler struct {
	transferService *services.TransferService
}

// NewTransferHandler creates a new transfer handler
func NewTransferHandler() *TransferHandler {
	return &TransferHandler{
		transferService: services.NewTransferService(),
	}
}

// InitiateTransferRequest represents the request to initiate a transfer
type InitiateTransferRequest struct {
	StipendID     string `json:"stipend_id"`
	PaymentMethod string `json:"payment_method"`
}

// TransferResponse represents the response for transfer operations
type TransferResponse struct {
	ID                 string `json:"id"`
	StipendID          string `json:"stipend_id"`
	StudentID          string `json:"student_id"`
	Amount             float64 `json:"amount"`
	Status             string `json:"status"`
	ReferenceNumber    string `json:"reference_number,omitempty"`
	ErrorMessage       string `json:"error_message,omitempty"`
	PaymentMethod      string `json:"payment_method"`
	DestinationAccount string `json:"destination_account"`
	DestinationBank    string `json:"destination_bank"`
	InitiatedAt        string `json:"initiated_at"`
	ProcessedAt        string `json:"processed_at,omitempty"`
	CompletedAt        string `json:"completed_at,omitempty"`
}

// InitiateTransfer initiates a money transfer for a stipend
// POST /api/transfers/initiate
func (th *TransferHandler) InitiateTransfer(w http.ResponseWriter, r *http.Request) {
	var req InitiateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.StipendID == "" {
		http.Error(w, "stipend_id is required", http.StatusBadRequest)
		return
	}

	if req.PaymentMethod == "" {
		req.PaymentMethod = "BANK_TRANSFER" // Default payment method
	}

	stipendID, err := uuid.Parse(req.StipendID)
	if err != nil {
		http.Error(w, "Invalid stipend_id format", http.StatusBadRequest)
		return
	}

	log.Printf("Initiating transfer for stipend: %s with method: %s", stipendID, req.PaymentMethod)

	transaction, err := th.transferService.InitiateTransfer(stipendID, req.PaymentMethod)
	if err != nil {
		log.Printf("Error initiating transfer: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := transactionToResponse(transaction)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ProcessTransfer processes a pending transfer
// POST /api/transfers/{transactionID}/process
func (th *TransferHandler) ProcessTransfer(w http.ResponseWriter, r *http.Request) {
	transactionIDStr := chi.URLParam(r, "transactionID")
	if transactionIDStr == "" {
		http.Error(w, "transactionID is required", http.StatusBadRequest)
		return
	}

	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		http.Error(w, "Invalid transactionID format", http.StatusBadRequest)
		return
	}

	log.Printf("Processing transfer: %s", transactionID)

	transaction, err := th.transferService.ProcessTransfer(transactionID)
	if err != nil {
		log.Printf("Error processing transfer: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := transactionToResponse(transaction)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTransferStatus retrieves the status of a transfer
// GET /api/transfers/{transactionID}/status
func (th *TransferHandler) GetTransferStatus(w http.ResponseWriter, r *http.Request) {
	transactionIDStr := chi.URLParam(r, "transactionID")
	if transactionIDStr == "" {
		http.Error(w, "transactionID is required", http.StatusBadRequest)
		return
	}

	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		http.Error(w, "Invalid transactionID format", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching transfer status: %s", transactionID)

	transaction, err := th.transferService.GetTransactionStatus(transactionID)
	if err != nil {
		log.Printf("Error fetching transfer status: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := transactionToResponse(transaction)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTransactionsByStipend retrieves all transactions for a stipend
// GET /api/stipends/{stipendID}/transactions
func (th *TransferHandler) GetTransactionsByStipend(w http.ResponseWriter, r *http.Request) {
	stipendIDStr := chi.URLParam(r, "stipendID")
	if stipendIDStr == "" {
		http.Error(w, "stipendID is required", http.StatusBadRequest)
		return
	}

	stipendID, err := uuid.Parse(stipendIDStr)
	if err != nil {
		http.Error(w, "Invalid stipendID format", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching transactions for stipend: %s", stipendID)

	transactions, err := th.transferService.GetTransactionsByStipend(stipendID)
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]TransferResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = *transactionToResponse(&tx)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// GetTransactionsByStudent retrieves all transactions for a student
// GET /api/students/{studentID}/transactions
func (th *TransferHandler) GetTransactionsByStudent(w http.ResponseWriter, r *http.Request) {
	studentIDStr := chi.URLParam(r, "studentID")
	if studentIDStr == "" {
		http.Error(w, "studentID is required", http.StatusBadRequest)
		return
	}

	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid studentID format", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching transactions for student: %s", studentID)

	transactions, err := th.transferService.GetTransactionsByStudent(studentID)
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]TransferResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = *transactionToResponse(&tx)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// CancelTransfer cancels a pending transfer
// POST /api/transfers/{transactionID}/cancel
func (th *TransferHandler) CancelTransfer(w http.ResponseWriter, r *http.Request) {
	transactionIDStr := chi.URLParam(r, "transactionID")
	if transactionIDStr == "" {
		http.Error(w, "transactionID is required", http.StatusBadRequest)
		return
	}

	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		http.Error(w, "Invalid transactionID format", http.StatusBadRequest)
		return
	}

	var cancelReq struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&cancelReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Cancelling transfer: %s with reason: %s", transactionID, cancelReq.Reason)

	transaction, err := th.transferService.CancelTransfer(transactionID, cancelReq.Reason)
	if err != nil {
		log.Printf("Error cancelling transfer: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := transactionToResponse(transaction)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RetryFailedTransfer retries a failed transfer
// POST /api/transfers/{transactionID}/retry
func (th *TransferHandler) RetryFailedTransfer(w http.ResponseWriter, r *http.Request) {
	transactionIDStr := chi.URLParam(r, "transactionID")
	if transactionIDStr == "" {
		http.Error(w, "transactionID is required", http.StatusBadRequest)
		return
	}

	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		http.Error(w, "Invalid transactionID format", http.StatusBadRequest)
		return
	}

	log.Printf("Retrying failed transfer: %s", transactionID)

	transaction, err := th.transferService.RetryFailedTransfer(transactionID)
	if err != nil {
		log.Printf("Error retrying transfer: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := transactionToResponse(transaction)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// transactionToResponse converts a Transaction model to a TransferResponse
func transactionToResponse(transaction *models.Transaction) *TransferResponse {
	if transaction == nil {
		return nil
	}

	refNum := ""
	if transaction.ReferenceNumber.Valid {
		refNum = transaction.ReferenceNumber.String
	}

	response := &TransferResponse{
		ID:                 transaction.ID.String(),
		StipendID:          transaction.StipendID.String(),
		StudentID:          transaction.StudentID.String(),
		Amount:             transaction.Amount,
		Status:             transaction.Status,
		ReferenceNumber:    refNum,
		ErrorMessage:       transaction.ErrorMessage,
		PaymentMethod:      transaction.PaymentMethod,
		DestinationAccount: transaction.DestinationAccount,
		DestinationBank:    transaction.DestinationBank,
		InitiatedAt:        transaction.InitiatedAt.String(),
	}

	if transaction.ProcessedAt != nil {
		response.ProcessedAt = transaction.ProcessedAt.String()
	}

	if transaction.CompletedAt != nil {
		response.CompletedAt = transaction.CompletedAt.String()
	}

	return response
}

// Stub Transaction type for transactionToResponse - the real one is in models
// This is a workaround to avoid circular imports
type Transaction struct {
	ID                 uuid.UUID
	StipendID          uuid.UUID
	StudentID          uuid.UUID
	Amount             float64
	Status             string
	ReferenceNumber    string
	ErrorMessage       string
	PaymentMethod      string
	DestinationAccount string
	DestinationBank    string
	InitiatedAt        interface{}
	ProcessedAt        interface{}
	CompletedAt        interface{}
}
