package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// BankingServiceClient handles communication with Banking Service
type BankingServiceClient struct {
	BaseURL string
	Client  *http.Client
}

// StudentBankDetails represents bank details to sync
type StudentBankDetails struct {
	StudentID         uint   `json:"student_id"`
	BankID            uint   `json:"bank_id"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
}

// NewBankingServiceClient creates a new banking service client
func NewBankingServiceClient() *BankingServiceClient {
	baseURL := os.Getenv("BANKING_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8083" // Default for local development
	}

	return &BankingServiceClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// SyncBankDetails sends student bank details to Banking Service
func (c *BankingServiceClient) SyncBankDetails(details StudentBankDetails) error {
	url := fmt.Sprintf("%s/api/student-bank-details", c.BaseURL)

	jsonData, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal bank details: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to banking service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("banking service returned status: %d", resp.StatusCode)
	}

	return nil
}

// UpdateBankDetails updates existing bank details in Banking Service
func (c *BankingServiceClient) UpdateBankDetails(id uint, details StudentBankDetails) error {
	url := fmt.Sprintf("%s/api/student-bank-details/%d", c.BaseURL, id)

	jsonData, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal bank details: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to banking service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("banking service returned status: %d", resp.StatusCode)
	}

	return nil
}
