package main

import (
	"api-gateway/internal/router"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.RequestID)
    
    // CORS middleware - must be before routes
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    })
    
    router.SetupRoutes(r)
    
    // Catch all - log what routes don't match
    r.NotFound(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("404 Not Found: %s %s", r.Method, r.RequestURI)
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "404 Not Found: %s %s\n", r.Method, r.RequestURI)
    })

    log.Println("API Gateway running on port 8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}