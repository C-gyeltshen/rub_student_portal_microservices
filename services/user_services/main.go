package main

import (
    "log"
    "net/http"
    "os"
    "user_services/database"
    "user_services/handlers"

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

    // Menu endpoints (note: no /api prefix)
    r.Get("/api/get/users", handlers.GetUsers)
    r.Post("/api/create/user", handlers.CreateUsers)
    r.Get("/api/get/user/{id}", handlers.GetuserById)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8082"
    }

    log.Printf("User Service starting on :%s", port)
    http.ListenAndServe(":"+port, r)
}