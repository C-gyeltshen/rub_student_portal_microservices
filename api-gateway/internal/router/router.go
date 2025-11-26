package router

import (
	"api-gateway/internal/handlers"
	authmw "api-gateway/internal/middleware"
	"api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(r *chi.Mux) {
	// Add basic middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	
	// Public routes (no authentication required)
	r.Get("/", handlers.LoginInstructions)
	r.Get("/dashboard", handlers.Dashboard)
	r.Get("/health", handlers.Health)
	r.Get("/login-instructions", handlers.LoginInstructions)
	
	// Protected routes that require authentication
	r.Group(func(r chi.Router) {
		r.Use(authmw.AuthMiddleware)
		
		// User profile - all authenticated users
		r.Get("/profile", handlers.UserProfile)
		
		// Admin only routes
		r.Group(func(r chi.Router) {
			r.Use(authmw.AuthorizeAdmin())
			r.Get("/admin/dashboard", handlers.AdminDashboard)
		})
		
		// Finance officer and admin routes
		r.Group(func(r chi.Router) {
			r.Use(authmw.AuthorizeFinanceOfficer())
			r.Get("/finance/dashboard", handlers.FinanceOfficerDashboard)
		})
		
		// Student, finance officer, and admin routes
		r.Group(func(r chi.Router) {
			r.Use(authmw.AuthorizeStudent())
			r.Get("/student/dashboard", handlers.StudentDashboard)
		})
		
		// Protected API routes that proxy to downstream services
		
		// User service routes (Admin and Finance Officer only)
		r.Group(func(r chi.Router) {
			r.Use(authmw.AuthorizeRoles(authmw.RoleAdmin, authmw.RoleFinanceOfficer))
			
			r.Get("/api/users", proxy.ForwardToUserService)
			r.Post("/api/users", proxy.ForwardToUserService)
			r.Get("/api/users/*", proxy.ForwardToUserService)
			r.Post("/api/users/*", proxy.ForwardToUserService)
			r.Patch("/api/users/*", proxy.ForwardToUserService)
			r.Delete("/api/users/*", proxy.ForwardToUserService)
		})
		
		// Banking service routes (Different authorization levels)
		r.Group(func(r chi.Router) {
			// Read-only access for all authenticated users
			r.Group(func(r chi.Router) {
				r.Use(authmw.AuthorizeStudent())
				r.Get("/api/banks", proxy.ForwardToBankingService)
				r.Get("/api/banks/*", proxy.ForwardToBankingService)
			})
			
			// Write access for Admin and Finance Officer only
			r.Group(func(r chi.Router) {
				r.Use(authmw.AuthorizeRoles(authmw.RoleAdmin, authmw.RoleFinanceOfficer))
				r.Post("/api/banks/*", proxy.ForwardToBankingService)
				r.Put("/api/banks/*", proxy.ForwardToBankingService)
				r.Delete("/api/banks/*", proxy.ForwardToBankingService)
			})
			
			// Student bank details - students can manage their own, others need higher permissions
			r.Group(func(r chi.Router) {
				r.Use(authmw.AuthorizeStudent())
				r.Get("/api/student-bank-details/*", proxy.ForwardToBankingService)
				r.Post("/api/student-bank-details", proxy.ForwardToBankingService)
				r.Put("/api/student-bank-details/*", proxy.ForwardToBankingService)
			})
			
			// Admin/Finance Officer can manage all student bank details
			r.Group(func(r chi.Router) {
				r.Use(authmw.AuthorizeRoles(authmw.RoleAdmin, authmw.RoleFinanceOfficer))
				r.Get("/api/student-bank-details", proxy.ForwardToBankingService)
				r.Delete("/api/student-bank-details/*", proxy.ForwardToBankingService)
			})
		})
	})
}