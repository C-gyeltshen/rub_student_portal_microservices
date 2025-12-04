package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// BankingServiceClient interacts with Banking Service
type BankingServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// Bank represents bank information from Banking Service
type Bank struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

// StudentBankDetails represents student's bank account information
type StudentBankDetails struct {
	ID                string `json:"id"`
	StudentID         string `json:"student_id"`
	BankID            string `json:"bank_id"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
	CreatedAt         string `json:"created_at"`
	ModifiedAt        string `json:"modified_at"`
}

// NewBankingServiceClient creates a new Banking Service client
func NewBankingServiceClient() *BankingServiceClient {
	baseURL := os.Getenv("BANKING_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://banking_services:8083" // Default service URL
	}

	return &BankingServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetBank fetches bank information by ID
func (c *BankingServiceClient) GetBank(bankID string) (*Bank, error) {
	url := fmt.Sprintf("%s/api/banks/%s", c.baseURL, bankID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bank: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("banking service returned %d: %s", resp.StatusCode, string(body))
	}

	var bank Bank
	if err := json.NewDecoder(resp.Body).Decode(&bank); err != nil {
		return nil, fmt.Errorf("failed to parse bank response: %w", err)
	}

	return &bank, nil
}

// GetStudentBankDetails fetches student's bank account information
func (c *BankingServiceClient) GetStudentBankDetails(studentID string) (*StudentBankDetails, error) {
	url := fmt.Sprintf("%s/api/student-bank-details/%s", c.baseURL, studentID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch student bank details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no bank details found for student %s", studentID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("banking service returned %d: %s", resp.StatusCode, string(body))
	}

	var bankDetails StudentBankDetails
	if err := json.NewDecoder(resp.Body).Decode(&bankDetails); err != nil {
		return nil, fmt.Errorf("failed to parse bank details response: %w", err)
	}

	return &bankDetails, nil
}

// ValidateStudentBankDetails checks if student has valid bank account for stipend
func (c *BankingServiceClient) ValidateStudentBankDetails(studentID string) error {
	bankDetails, err := c.GetStudentBankDetails(studentID)
	if err != nil {
		return err
	}

	if bankDetails == nil {
		return fmt.Errorf("bank details not found for student")
	}

	if bankDetails.AccountNumber == "" || bankDetails.AccountHolderName == "" {
		return fmt.Errorf("incomplete bank details for student")
	}

	// Verify bank exists
	bank, err := c.GetBank(bankDetails.BankID)
	if err != nil {
		return fmt.Errorf("failed to verify bank: %w", err)
	}

	if bank == nil {
		return fmt.Errorf("associated bank not found")
	}

	return nil
}

// GetAllBanks fetches list of all banks
func (c *BankingServiceClient) GetAllBanks() ([]*Bank, error) {
	url := fmt.Sprintf("%s/api/banks", c.baseURL)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch banks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("banking service returned %d: %s", resp.StatusCode, string(body))
	}

	var banks []*Bank
	if err := json.NewDecoder(resp.Body).Decode(&banks); err != nil {
		return nil, fmt.Errorf("failed to parse banks response: %w", err)
	}

	return banks, nil
}
