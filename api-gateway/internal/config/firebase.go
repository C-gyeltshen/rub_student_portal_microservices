package config

import (
	"context"
	"encoding/json"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	FirebaseApp    *firebase.App
	FirebaseAuth   *auth.Client
	FirebaseConfig *Config
)

type Config struct {
	ProjectID    string `json:"project_id"`
	ClientEmail  string `json:"client_email"`
	PrivateKey   string `json:"private_key"`
	DatabaseURL  string `json:"database_url,omitempty"`
}

// InitializeFirebase initializes Firebase Admin SDK
func InitializeFirebase() error {
	// Load environment variables
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	clientEmail := os.Getenv("FIREBASE_CLIENT_EMAIL")
	privateKey := os.Getenv("FIREBASE_PRIVATE_KEY")

	if projectID == "" || clientEmail == "" || privateKey == "" {
		log.Fatal("Firebase environment variables are required: FIREBASE_PROJECT_ID, FIREBASE_CLIENT_EMAIL, FIREBASE_PRIVATE_KEY")
	}

	// Create Firebase config
	FirebaseConfig = &Config{
		ProjectID:   projectID,
		ClientEmail: clientEmail,
		PrivateKey:  privateKey,
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
		return err
	}

	// Initialize Firebase app
	opt := option.WithCredentialsJSON(credentialsBytes)
	config := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return err
	}

	FirebaseApp = app

	// Initialize Auth client
	authClient, err := app.Auth(context.Background())
	if err != nil {
		return err
	}

	FirebaseAuth = authClient

	log.Println("Firebase Admin SDK initialized successfully")
	return nil
}

// GetAuthClient returns the Firebase Auth client
func GetAuthClient() *auth.Client {
	return FirebaseAuth
}

// GetProjectID returns the Firebase project ID
func GetProjectID() string {
	return FirebaseConfig.ProjectID
}