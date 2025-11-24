package main

import (
	"log"
	"net/http"
	"os"
	"banking_services/database"
	"banking_services/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Bank endpoints
	r.Get("/api/banks", handlers.GetBanks)
	r.Post("/api/banks", handlers.CreateBank)
	r.Get("/api/banks/{id}", handlers.GetBankById)
	r.Put("/api/banks/{id}", handlers.UpdateBank)
	r.Delete("/api/banks/{id}", handlers.DeleteBank)

	// Student Bank Details endpoints
	r.Get("/api/student-bank-details", handlers.GetStudentBankDetails)
	r.Post("/api/student-bank-details", handlers.CreateStudentBankDetails)
	r.Get("/api/student-bank-details/{id}", handlers.GetStudentBankDetailsById)
	r.Get("/api/student-bank-details/student/{studentId}", handlers.GetStudentBankDetailsByStudentId)
	r.Put("/api/student-bank-details/{id}", handlers.UpdateStudentBankDetails)
	r.Delete("/api/student-bank-details/{id}", handlers.DeleteStudentBankDetails)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Banking Service starting on :%s", port)
	http.ListenAndServe(":"+port, r)
}
