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
		grpcPort = "50052"
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
	log.Println("DEBUG: Initializing stipend handler...")
	stipendHandler := handlers.NewStipendHandler()
	log.Println("DEBUG: Stipend handler initialized successfully")
	
	log.Println("DEBUG: Initializing deduction handler...")
	deductionHandler := handlers.NewDeductionHandler()
	log.Println("DEBUG: Deduction handler initialized successfully")

	log.Println("DEBUG: Initializing transfer handler...")
	transferHandler := handlers.NewTransferHandler()
	log.Println("DEBUG: Transfer handler initialized successfully")

	log.Println("DEBUG: Initializing search handler...")
	searchHandler := handlers.NewSearchHandler()
	log.Println("DEBUG: Search handler initialized successfully")

	log.Println("DEBUG: Initializing report handler...")
	reportHandler := handlers.NewReportHandler()
	log.Println("DEBUG: Report handler initialized successfully")

	log.Println("DEBUG: Initializing audit handler...")
	auditHandler := handlers.NewAuditHandler()
	log.Println("DEBUG: Audit handler initialized successfully")

	// API Routes
	log.Println("DEBUG: Registering API routes...")
	r.Route("/api", func(r chi.Router) {
		// Stipend endpoints
		r.Post("/stipends", stipendHandler.CreateStipend)
		r.Post("/stipends/calculate", stipendHandler.CalculateStipendWithDeductions)
		r.Post("/stipends/calculate/monthly", stipendHandler.CalculateMonthlyStipend)
		r.Post("/stipends/calculate/annual", stipendHandler.CalculateAnnualStipend)
		r.Get("/stipends/{stipendID}", stipendHandler.GetStipend)
		r.Get("/stipends/{stipendID}/deductions", stipendHandler.GetStipendDeductions)
		r.Get("/stipends/{stipendID}/transactions", transferHandler.GetTransactionsByStipend)
		r.Patch("/stipends/{stipendID}/payment-status", stipendHandler.UpdateStipendPaymentStatus)

		// Student stipends endpoint
		r.Get("/students/{studentID}/stipends", stipendHandler.GetStudentStipends)
		r.Get("/students/{studentID}/transactions", transferHandler.GetTransactionsByStudent)

		// Deduction rule endpoints
		r.Post("/deduction-rules", deductionHandler.CreateDeductionRule)
		r.Get("/deduction-rules", deductionHandler.ListDeductionRules)
		r.Get("/deduction-rules/{ruleID}", deductionHandler.GetDeductionRule)

		// Money Transfer endpoints
		r.Post("/transfers/initiate", transferHandler.InitiateTransfer)
		r.Post("/transfers/{transactionID}/process", transferHandler.ProcessTransfer)
		r.Get("/transfers/{transactionID}/status", transferHandler.GetTransferStatus)
		r.Post("/transfers/{transactionID}/cancel", transferHandler.CancelTransfer)
		r.Post("/transfers/{transactionID}/retry", transferHandler.RetryFailedTransfer)

		// Search endpoints
		r.Get("/search/stipends", searchHandler.SearchStipends)
		r.Get("/search/deduction-rules", searchHandler.SearchDeductionRules)
		r.Get("/search/transactions", searchHandler.SearchTransactions)

		// Report endpoints
		r.Get("/reports/disbursement", reportHandler.GetDisbursementReport)
		r.Get("/reports/deductions", reportHandler.GetDeductionReport)
		r.Get("/reports/transactions", reportHandler.GetTransactionReport)
		r.Get("/reports/export/stipends", reportHandler.ExportStipendsCsv)
		r.Get("/reports/export/deductions", reportHandler.ExportDeductionsCsv)
		r.Get("/reports/export/transactions", reportHandler.ExportTransactionsCsv)
		r.Get("/reports/export/pdf/stipends", reportHandler.ExportStipendsPdf)
		r.Get("/reports/export/pdf/deductions", reportHandler.ExportDeductionsPdf)
		r.Get("/reports/export/pdf/transactions", reportHandler.ExportTransactionsPdf)

		// Audit log endpoints
		r.Get("/audit-logs", auditHandler.GetAuditLogs)
		r.Get("/audit-logs/{entity_type}/{entity_id}", auditHandler.GetAuditLogsByEntity)
		r.Get("/audit-logs/officer/{officer}", auditHandler.GetAuditLogsByOfficer)
	})
	log.Println("DEBUG: API routes registered successfully")

	log.Printf("REST Server starting on port %s", port)
	return http.ListenAndServe(":"+port, r)
}

