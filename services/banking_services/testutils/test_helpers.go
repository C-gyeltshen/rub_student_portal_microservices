package testutils

import (
	"banking_services/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreateTestBank returns a test bank object
func CreateTestBank() *models.Bank {
	return &models.Bank{
		ID:   uuid.New(),
		Name: "Test Bank",
	}
}

// CreateTestStudentBankDetails returns test student bank details
func CreateTestStudentBankDetails() *models.StudentBankDetails {
	return &models.StudentBankDetails{
		ID:                uuid.New(),
		StudentID:         uuid.New(),
		BankID:            uuid.New(),
		AccountNumber:     "1234567890",
		AccountHolderName: "John Doe",
	}
}

// SetupTestRouter creates a chi router for testing
func SetupTestRouter() *chi.Mux {
	return chi.NewRouter()
}

// MakeRequest is a helper to make HTTP requests in tests
func MakeRequest(method, url string, body interface{}) (*http.Request, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}
	return http.NewRequest(method, url, reqBody)
}

// ExecuteRequest executes a request and returns response recorder
func ExecuteRequest(req *http.Request, router *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
