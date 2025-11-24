package main

import (
	"finance_service/database"
	"finance_service/handlers"
	"finance_service/internal/grpc"
	pb "finance_service/pkg/pb"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	grpclib "google.golang.org/grpc"
)

func main() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize database (create indexes, constraints, seed data)
	if err := database.InitializeFinanceDatabase(); err != nil {
		log.Printf("Warning: Failed to initialize finance database: %v", err)
	}

	// Get ports from environment
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8084"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	// Start gRPC and REST servers in parallel
	var wg sync.WaitGroup
	var grpcErr error
	wg.Add(2)

	// Start gRPC server
	go func() {
		defer wg.Done()
		if err := startGRPCServer(grpcPort); err != nil {
			grpcErr = err
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Start REST server
	go func() {
		defer wg.Done()
		if err := startRESTServer(httpPort); err != nil {
			log.Printf("REST server error: %v (continuing without REST server)", err)
		}
	}()

	wg.Wait()

	// If gRPC server failed, exit with error
	if grpcErr != nil {
		log.Fatalf("Fatal error: gRPC server failed: %v", grpcErr)
	}
}

// startGRPCServer starts the gRPC server
func startGRPCServer(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Printf("gRPC Server listening on port %s", port)

	server := grpclib.NewServer()

	// Register services
	pb.RegisterStipendServiceServer(server, grpc.NewStipendServiceServer())
	pb.RegisterDeductionServiceServer(server, grpc.NewDeductionServiceServer())

	if err := server.Serve(listener); err != nil {
		return err
	}

	return nil
}

// startRESTServer starts the REST API server
func startRESTServer(port string) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Finance Service is healthy"))
	})

	// Initialize handlers
	stipendHandler := handlers.NewStipendHandler()
	deductionHandler := handlers.NewDeductionHandler()

	// API Routes
	r.Route("/api", func(r chi.Router) {
		// Stipend endpoints
		r.Post("/stipends", stipendHandler.CreateStipend)
		r.Post("/stipends/calculate", stipendHandler.CalculateStipendWithDeductions)
		r.Post("/stipends/calculate/monthly", stipendHandler.CalculateMonthlyStipend)
		r.Post("/stipends/calculate/annual", stipendHandler.CalculateAnnualStipend)
		r.Get("/stipends/{stipendID}", stipendHandler.GetStipend)
		r.Get("/stipends/{stipendID}/deductions", stipendHandler.GetStipendDeductions)
		r.Patch("/stipends/{stipendID}/payment-status", stipendHandler.UpdateStipendPaymentStatus)

		// Student stipends endpoint
		r.Get("/students/{studentID}/stipends", stipendHandler.GetStudentStipends)

		// Deduction rule endpoints
		r.Post("/deduction-rules", deductionHandler.CreateDeductionRule)
		r.Get("/deduction-rules", deductionHandler.ListDeductionRules)
		r.Get("/deduction-rules/{ruleID}", deductionHandler.GetDeductionRule)
	})

	log.Printf("REST Server starting on port %s", port)
	return http.ListenAndServe(":"+port, r)
}

