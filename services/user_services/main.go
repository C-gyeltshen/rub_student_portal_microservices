package main

import (
	"log"
	"net/http"
	"os"
	"user_services/database"
	"user_services/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Menu endpoints (note: no /api prefix)
	r.Get("/users", handlers.GetUsers)
	r.Post("/users/create", handlers.CreateUsers)
	r.Post("/users/create/finance-officer", handlers.CreateFinanceOfficer)
	r.Get("/users/{id}", handlers.GetuserById)
	r.Get("/users/role/{roleId}", handlers.GetUsersByRoleId)
	r.Delete("/users/role/{roleId}", handlers.DeleteUsersByRoleId)

	r.Post("/users/create/roles", handlers.CreateRole)
	r.Get("/users/get/roles", handlers.GetRoles)
	r.Get("/users/get/role/{id}", handlers.GetRoleById)
	r.Patch("/users/update/role/{id}", handlers.UpdateRole)
	r.Delete("/users/delete/role/{id}", handlers.DeleteRole)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("User Service starting on :%s", port)
	http.ListenAndServe(":"+port, r)
}
