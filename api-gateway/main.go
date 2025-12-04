package main

import (
	gwMiddleware "api-gateway/internal/middleware"
	"api-gateway/internal/router"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    
    // Add middleware
    r.Use(middleware.Logger)
    r.Use(gwMiddleware.CORSMiddleware())
    
    router.SetupRoutes(r)

    log.Println("API Gateway running on port 8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}