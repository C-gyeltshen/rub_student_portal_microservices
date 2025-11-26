package main

import (
	"log"
	"net/http"

	"api-gateway/internal/config"
	"api-gateway/internal/router"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading .env file:", err)
	}

	// Initialize Firebase Admin SDK
	if err := config.InitializeFirebase(); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	r := chi.NewRouter()
	router.SetupRoutes(r)

	log.Println("API Gateway running on port 8080")
	log.Println("Firebase authentication enabled")
	log.Println("Public endpoints: /, /dashboard, /health, /login-instructions")
	log.Println("Protected endpoints require Firebase ID token in Authorization header")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}