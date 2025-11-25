package router

import (
	"api-gateway/internal/proxy"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux){
	// User service routes
    r.Get("/api/users", proxy.ForwardToBankingService)
    r.Post("/api/users", proxy.ForwardToBankingService)
    r.Get("/api/users/*", proxy.ForwardToBankingService)
    r.Post("/api/users/*", proxy.ForwardToBankingService)      // Add this line
    r.Patch("/api/users/*", proxy.ForwardToBankingService)     // Add this line
    r.Delete("/api/users/*", proxy.ForwardToBankingService)
	
	// Banking service routes
    // r.Post("/api/banks/*", proxy.ForwardToBankingService)
    r.Get("/api/banks", proxy.ForwardToBankingService)
    r.Post("/api/banks", proxy.ForwardToBankingService)
    r.Get("/api/banks/*", proxy.ForwardToBankingService)
    r.Post("/api/banks/*", proxy.ForwardToBankingService)      // Add this line
    r.Patch("/api/banks/*", proxy.ForwardToBankingService)     // Add this line
    r.Delete("/api/banks/*", proxy.ForwardToBankingService)
}