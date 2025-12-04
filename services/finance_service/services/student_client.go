package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// StudentServiceClient interacts with Student Service
type StudentServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// Student represents student information from Student Service
type Student struct {
	ID              string `json:"id"`
	UserID          string `json:"user_id"`
	Name            string `json:"name"`
	RubIDCardNumber string `json:"rub_id_card_number"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	DateOfBirth     string `json:"date_of_birth"`
	CollegeID       string `json:"college_id"`
	ProgramID       string `json:"program_id"`
	CreatedAt       string `json:"created_at"`
	ModifiedAt      string `json:"modified_at"`
}

// NewStudentServiceClient creates a new Student Service client
func NewStudentServiceClient() *StudentServiceClient {
	baseURL := os.Getenv("STUDENT_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://user_services:8082" // Default service URL
	}

	return &StudentServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetStudent fetches student information by ID
func (c *StudentServiceClient) GetStudent(studentID string) (*Student, error) {
	url := fmt.Sprintf("%s/api/students/%s", c.baseURL, studentID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch student: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("student service returned %d: %s", resp.StatusCode, string(body))
	}

	var student Student
	if err := json.NewDecoder(resp.Body).Decode(&student); err != nil {
		return nil, fmt.Errorf("failed to parse student response: %w", err)
	}

	return &student, nil
}

// ValidateStudent checks if student exists and is eligible for stipend
func (c *StudentServiceClient) ValidateStudent(studentID string) error {
	student, err := c.GetStudent(studentID)
	if err != nil {
		return err
	}

	if student == nil {
		return fmt.Errorf("student not found")
	}

	if student.Email == "" || student.Name == "" {
		return fmt.Errorf("student record is incomplete")
	}

	return nil
}

// GetStudentByCardNumber fetches student by RUB ID card number
func (c *StudentServiceClient) GetStudentByCardNumber(cardNumber string) (*Student, error) {
	url := fmt.Sprintf("%s/api/students/card/%s", c.baseURL, cardNumber)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch student by card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("student with card number %s not found", cardNumber)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("student service returned %d: %s", resp.StatusCode, string(body))
	}

	var student Student
	if err := json.NewDecoder(resp.Body).Decode(&student); err != nil {
		return nil, fmt.Errorf("failed to parse student response: %w", err)
	}

	return &student, nil
}
