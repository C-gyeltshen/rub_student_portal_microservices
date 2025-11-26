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
	r.Get("/banks", handlers.GetBanks)
	r.Post("/banks", handlers.CreateBank)
	r.Get("/banks/{id}", handlers.GetBankById)
	r.Patch("/banks/{id}", handlers.UpdateBank)
	r.Delete("/banks/{id}", handlers.DeleteBank)

	// Student Bank Details endpoints
	r.Get("/banks/get/student-bank-details", handlers.GetStudentBankDetails)
	r.Post("/banks/create/student-bank-details", handlers.CreateStudentBankDetails)
	r.Get("/banks/get/student-bank-details/{id}", handlers.GetStudentBankDetailsById)
	r.Get("/banks/get/student-bank-details/student/{studentId}", handlers.GetStudentBankDetailsByStudentId)
	r.Patch("/banks/update/student-bank-details/{id}", handlers.UpdateStudentBankDetails)
	r.Delete("/banks/delete/student-bank-details/{id}", handlers.DeleteStudentBankDetails)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Banking Service starting on :%s", port)
	http.ListenAndServe(":"+port, r)
}