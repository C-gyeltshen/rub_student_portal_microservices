package main
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run set_admin_role.go <user_uid> <role>")
		log.Fatal("Example: go run set_admin_role.go dWvtrkWzwfaeRFiW6ARnwtg6ErI2 admin")
	}
	
	userUID := os.Args[1]
	role := os.Args[2]
	
	// Validate role
	validRoles := []string{"admin", "finance_officer", "student"}
	isValid := false
	for _, validRole := range validRoles {
		if role == validRole {
			isValid = true
			break
		}
	}
	
	if !isValid {
		log.Fatalf("Invalid role: %s. Valid roles: admin, finance_officer, student", role)
	}

	// Load .env file from parent directory
	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("Warning: No .env file found or error loading .env file:", err)
	}

	// Get environment variables
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	clientEmail := os.Getenv("FIREBASE_CLIENT_EMAIL")
	privateKey := os.Getenv("FIREBASE_PRIVATE_KEY")

	if projectID == "" || clientEmail == "" || privateKey == "" {
		log.Fatal("Firebase environment variables are required: FIREBASE_PROJECT_ID, FIREBASE_CLIENT_EMAIL, FIREBASE_PRIVATE_KEY")
	}

	// Create service account credentials JSON
	credentialsJSON := map[string]interface{}{
		"type":                        "service_account",
		"project_id":                  projectID,
		"private_key_id":              "",
		"private_key":                 privateKey,
		"client_email":                clientEmail,
		"client_id":                   "",
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/" + clientEmail,
	}

	credentialsBytes, err := json.Marshal(credentialsJSON)
	if err != nil {
		log.Fatalf("Error creating credentials: %v", err)
	}

	// Initialize Firebase app
	opt := option.WithCredentialsJSON(credentialsBytes)
	config := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error getting Auth client: %v", err)
	}

	// Verify user exists
	user, err := authClient.GetUser(context.Background(), userUID)
	if err != nil {
		log.Fatalf("Error verifying user %s: %v", userUID, err)
	}

	// Set admin custom claims
	claims := map[string]interface{}{
		"role": role,
	}

	err = authClient.SetCustomUserClaims(context.Background(), userUID, claims)
	if err != nil {
		log.Fatalf("Error setting custom claims: %v", err)
	}

	log.Printf("Successfully set '%s' role for user: %s (email: %s)", role, userUID, user.Email)
	log.Printf("The user will need to refresh their ID token to see the new role.")
}