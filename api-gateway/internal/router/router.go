package router

import (
	"api-gateway/internal/proxy"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux){
	// User service routes
    r.Get("/api/users", proxy.ForwardToUserService)
    r.Post("/api/users", proxy.ForwardToUserService)
    r.Get("/api/users/*", proxy.ForwardToUserService)
    r.Post("/api/users/*", proxy.ForwardToUserService)      // Add this line
    r.Patch("/api/users/*", proxy.ForwardToUserService)     // Add this line
    r.Delete("/api/users/*", proxy.ForwardToUserService)
	
	// Banking service routes
    // r.Post("/api/banks/*", proxy.ForwardToBankingService)
    r.Get("/api/banks", proxy.ForwardToBankingService)
    r.Post("/api/banks", proxy.ForwardToBankingService)
    r.Get("/api/banks/*", proxy.ForwardToBankingService)
    r.Post("/api/banks/*", proxy.ForwardToBankingService)      // Add this line
    r.Patch("/api/banks/*", proxy.ForwardToBankingService)     // Add this line
    r.Delete("/api/banks/*", proxy.ForwardToBankingService)

	// Student Management service routes
    r.Get("/api/students", proxy.ForwardToStudentService)
    r.Post("/api/students", proxy.ForwardToStudentService)
    r.Post("/api/students/bulk", proxy.ForwardToStudentService)
    r.Get("/api/students/*", proxy.ForwardToStudentService)
    r.Put("/api/students/*", proxy.ForwardToStudentService)
    r.Delete("/api/students/*", proxy.ForwardToStudentService)
    
    r.Get("/api/programs", proxy.ForwardToStudentService)
    r.Post("/api/programs", proxy.ForwardToStudentService)
    r.Get("/api/programs/*", proxy.ForwardToStudentService)
    r.Put("/api/programs/*", proxy.ForwardToStudentService)
    
    r.Get("/api/colleges", proxy.ForwardToStudentService)
    r.Post("/api/colleges", proxy.ForwardToStudentService)
    r.Get("/api/colleges/*", proxy.ForwardToStudentService)
    r.Put("/api/colleges/*", proxy.ForwardToStudentService)
    
    r.Get("/api/stipend/*", proxy.ForwardToStudentService)
    r.Post("/api/stipend/*", proxy.ForwardToStudentService)
    r.Put("/api/stipend/*", proxy.ForwardToStudentService)
    
    r.Get("/api/reports/*", proxy.ForwardToStudentService)
}