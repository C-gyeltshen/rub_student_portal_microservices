package router

import (
	"api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux){
	// Mount subrouters for each service
	// These use Mount which matches prefix
	
	// Finance service routes
	financeRouter := chi.NewRouter()
	financeRouter.Get("/*", proxy.ForwardToFinanceService)
	financeRouter.Post("/*", proxy.ForwardToFinanceService)
	r.Mount("/api/finance", financeRouter)
	
	// User service routes
	userRouter := chi.NewRouter()
	userRouter.Get("/*", proxy.ForwardToUserService)
	userRouter.Post("/*", proxy.ForwardToUserService)
	r.Mount("/api/users", userRouter)
	
	// Banking service routes
	bankRouter := chi.NewRouter()
	bankRouter.Get("/*", proxy.ForwardToBankingService)
	bankRouter.Post("/*", proxy.ForwardToBankingService)
	r.Mount("/api/banks", bankRouter)
}

