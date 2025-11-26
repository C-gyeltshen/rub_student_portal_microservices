package router

import (
	"api-gateway/internal/proxy"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux){
        log.Println("SetupRoutes called - registering routes")

        // Health check endpoint
        r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                w.Write([]byte(`{"status":"ok"}`))
        })

        // Test route for debugging
        r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                w.Write([]byte(`{"message":"test route works"}`))
        })

        log.Println("Routes registered successfully")

        // Student Management Service routes
        r.Route("/api/students", func(r chi.Router) {
                r.Get("/", proxy.ForwardToStudentService)
                r.Post("/", proxy.ForwardToStudentService)
                r.Get("/*", proxy.ForwardToStudentService)
                r.Post("/*", proxy.ForwardToStudentService)
                r.Put("/*", proxy.ForwardToStudentService)
                r.Delete("/*", proxy.ForwardToStudentService)
        })

        // Colleges routes (from Student Service)
        r.Route("/api/colleges", func(r chi.Router) {
                r.Get("/", proxy.ForwardToStudentService)
                r.Post("/", proxy.ForwardToStudentService)
                r.Get("/*", proxy.ForwardToStudentService)
                r.Post("/*", proxy.ForwardToStudentService)
                r.Put("/*", proxy.ForwardToStudentService)
                r.Delete("/*", proxy.ForwardToStudentService)
        })

        // Programs routes (from Student Service)
        r.Route("/api/programs", func(r chi.Router) {
                r.Get("/", proxy.ForwardToStudentService)
                r.Post("/", proxy.ForwardToStudentService)
                r.Get("/*", proxy.ForwardToStudentService)
                r.Post("/*", proxy.ForwardToStudentService)
                r.Put("/*", proxy.ForwardToStudentService)
                r.Delete("/*", proxy.ForwardToStudentService)
        })

        // Banking service routes
        r.Route("/api/banks", func(r chi.Router) {
                r.Get("/", proxy.ForwardToBankingService)
                r.Post("/", proxy.ForwardToBankingService)
                r.Get("/*", proxy.ForwardToBankingService)
                r.Post("/*", proxy.ForwardToBankingService)
                r.Put("/*", proxy.ForwardToBankingService)
                r.Delete("/*", proxy.ForwardToBankingService)
        })

        // User service routes
        r.Route("/api/users", func(r chi.Router) {
                r.Get("/", proxy.ForwardToUserService)
                r.Post("/", proxy.ForwardToUserService)
                r.Get("/*", proxy.ForwardToUserService)
                r.Post("/*", proxy.ForwardToUserService)
                r.Put("/*", proxy.ForwardToUserService)
                r.Delete("/*", proxy.ForwardToUserService)
        })

        // Finance service routes
        r.Route("/api/finance", func(r chi.Router) {
                r.Get("/", proxy.ForwardToFinanceService)
                r.Post("/", proxy.ForwardToFinanceService)
                r.Get("/*", proxy.ForwardToFinanceService)
                r.Post("/*", proxy.ForwardToFinanceService)
                r.Put("/*", proxy.ForwardToFinanceService)
                r.Delete("/*", proxy.ForwardToFinanceService)
        })
}
