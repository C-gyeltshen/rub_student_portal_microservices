package main

import (
    "log"
    "net/http"
    "os"
    "monolith/database"
    "monolith/handlers"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    // Connect to database
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "host=localhost user=postgres password=postgres dbname=student_cafe port=5432 sslmode=disable"
    }

    if err := database.Connect(dsn); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Setup router
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // User routes
    r.Post("/api/users", handlers.CreateUser)
    r.Get("/api/users/{id}", handlers.GetUser)

    //College
    r.Get("/api/college/{id}", handlers.GetCollege)
    r.Post("/api/colleges", handlers.CreateCollege)

    // Bank
    r.Get("/api/bank/{id}", handlers.GetBank)
    r.Post("/api/banks", handlers.CreateBank)

    // Program
    r.Get("/api/program/{id}", handlers.GetProgram)
    r.Post("/api/program", handlers.CreateProgram)

    // Student
    r.Get("/api/student/{id}", handlers.GetStudent)
    r.Post("/api/student", handlers.CreateStudent)

    log.Println("Monolith server starting on :8080")
    http.ListenAndServe(":8080", r)
}