package testutils

import (
	"banking_services/models"
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
)

// CreateTestBank returns a test bank object
func CreateTestBank() *models.Bank {
	return &models.Bank{
		Model: gorm.Model{ID: 1},
		Name:  "Test Bank",
	}
}

// CreateTestStudentBankDetails returns test student bank details
func CreateTestStudentBankDetails() *models.StudentBankDetails {
	return &models.StudentBankDetails{
		Model:             gorm.Model{ID: 1},
		StudentID:         123,
		BankID:            1,
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
