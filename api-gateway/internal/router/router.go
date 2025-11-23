package router

import (
	"api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux){
	r.Get("/api/users/*", proxy.ForwardToUserService)
    r.Post("/api/users/*", proxy.ForwardToUserService)
}