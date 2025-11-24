package router

import (
	"api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux){
	// User service routes
    r.Get("/api/users", proxy.ForwardToUserService)
    r.Post("/api/users", proxy.ForwardToUserService)
    r.Get("/api/users/*", proxy.ForwardToUserService)
    r.Post("/api/users/*", proxy.ForwardToUserService)
    r.Patch("/api/users/*", proxy.ForwardToUserService)
    r.Delete("/api/users/*", proxy.ForwardToUserService)
	
	// Banking service routes
    r.Get("/api/banks", proxy.ForwardToBankingService)
    r.Post("/api/banks", proxy.ForwardToBankingService)
    r.Get("/api/banks/*", proxy.ForwardToBankingService)
    r.Put("/api/banks/*", proxy.ForwardToBankingService)
    r.Delete("/api/banks/*", proxy.ForwardToBankingService)

	// Finance service routes - Main endpoint and all sub-paths
    r.Get("/api/finance", proxy.ForwardToFinanceService)
    r.Post("/api/finance", proxy.ForwardToFinanceService)
    r.Get("/api/finance/*", proxy.ForwardToFinanceService)
    r.Post("/api/finance/*", proxy.ForwardToFinanceService)
    r.Put("/api/finance/*", proxy.ForwardToFinanceService)
    r.Delete("/api/finance/*", proxy.ForwardToFinanceService)
    r.Patch("/api/finance/*", proxy.ForwardToFinanceService)
}

