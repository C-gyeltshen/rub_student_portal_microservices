package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// UserServiceClient interacts with User Service
type UserServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// User represents user information from User Service
type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	RoleID       string `json:"role_id"`
	CreatedAt    string `json:"created_at"`
	ModifiedAt   string `json:"modified_at"`
}

// Role represents user role information
type Role struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

// NewUserServiceClient creates a new User Service client
func NewUserServiceClient() *UserServiceClient {
	baseURL := os.Getenv("USER_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://user_services:8082" // Default service URL
	}

	return &UserServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetUser fetches user information by ID
func (c *UserServiceClient) GetUser(userID string) (*User, error) {
	url := fmt.Sprintf("%s/api/users/%s", c.baseURL, userID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user service returned %d: %s", resp.StatusCode, string(body))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	return &user, nil
}

// GetRole fetches role information by ID
func (c *UserServiceClient) GetRole(roleID string) (*Role, error) {
	url := fmt.Sprintf("%s/api/roles/%s", c.baseURL, roleID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch role: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user service returned %d: %s", resp.StatusCode, string(body))
	}

	var role Role
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, fmt.Errorf("failed to parse role response: %w", err)
	}

	return &role, nil
}

// ValidateUserPermission checks if user has finance-related role
func (c *UserServiceClient) ValidateUserPermission(userID string) (string, error) {
	user, err := c.GetUser(userID)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", fmt.Errorf("user not found")
	}

	role, err := c.GetRole(user.RoleID)
	if err != nil {
		return "", err
	}

	if role == nil {
		return "", fmt.Errorf("role not found")
	}

	// Check if user has appropriate role for finance operations
	// Allowed roles: admin, finance_officer, finance_manager
	switch role.Name {
	case "admin", "finance_officer", "finance_manager":
		return role.Name, nil
	default:
		return "", fmt.Errorf("user role '%s' not authorized for finance operations", role.Name)
	}
}

// GetUserRole fetches the user's role name
func (c *UserServiceClient) GetUserRole(userID string) (string, error) {
	user, err := c.GetUser(userID)
	if err != nil {
		return "", err
	}

	role, err := c.GetRole(user.RoleID)
	if err != nil {
		return "", err
	}

	return role.Name, nil
}

// ValidateUserExists checks if user exists and is active
func (c *UserServiceClient) ValidateUserExists(userID string) error {
	user, err := c.GetUser(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return fmt.Errorf("user not found")
	}

	if user.Email == "" {
		return fmt.Errorf("user record is incomplete")
	}

	return nil
}
