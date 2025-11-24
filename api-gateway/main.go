package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
    "api-gateway/internal/router"
)

func main() {
    r := chi.NewRouter()
    router.SetupRoutes(r)

    log.Println("API Gateway running on port 8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}