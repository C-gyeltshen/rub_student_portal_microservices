package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"student_management_service/database"
	grpcserver "student_management_service/grpc/server"
	"student_management_service/handlers"
	pb "student_management_service/pb/student"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Start gRPC server in a separate goroutine
	go startGRPCServer()

	// Start HTTP REST server
	startHTTPServer()
}

// startGRPCServer starts the gRPC server on port 50054
func startGRPCServer() {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50054"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStudentServiceServer(grpcServer, &grpcserver.StudentServer{})

	log.Printf("gRPC server listening on :%s", grpcPort)
	log.Println("gRPC services:")
	log.Println("  - StudentService (GetStudent, CreateStudent, UpdateStudent, etc.)")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

// startHTTPServer starts the HTTP REST server on port 8084
func startHTTPServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// ==================== Student Endpoints ====================
	r.Get("/api/students", handlers.GetStudents)
	r.Post("/api/students", handlers.CreateStudent)
	r.Get("/api/students/search", handlers.SearchStudents)
	r.Get("/api/students/{id}", handlers.GetStudentById)
	r.Get("/api/students/student-id/{studentId}", handlers.GetStudentByStudentId)
	r.Get("/api/students/program/{programId}", handlers.GetStudentsByProgram)
	r.Get("/api/students/college/{collegeId}", handlers.GetStudentsByCollege)
	r.Get("/api/students/status/{status}", handlers.GetStudentsByStatus)
	r.Put("/api/students/{id}", handlers.UpdateStudent)
	r.Delete("/api/students/{id}", handlers.DeleteStudent)

	// ==================== Program Endpoints ====================
	r.Get("/api/programs", handlers.GetPrograms)
	r.Post("/api/programs", handlers.CreateProgram)
	r.Get("/api/programs/{id}", handlers.GetProgramById)
	r.Put("/api/programs/{id}", handlers.UpdateProgram)

	// ==================== College Endpoints ====================
	r.Get("/api/colleges", handlers.GetColleges)
	r.Post("/api/colleges", handlers.CreateCollege)
	r.Get("/api/colleges/{id}", handlers.GetCollegeById)
	r.Put("/api/colleges/{id}", handlers.UpdateCollege)

	// ==================== Stipend Eligibility Endpoints ====================
	r.Get("/api/stipend/eligibility/{studentId}", handlers.CheckStipendEligibility)

	// ==================== Stipend Allocation Endpoints ====================
	r.Get("/api/stipend/allocations", handlers.GetStipendAllocations)
	r.Post("/api/stipend/allocations", handlers.CreateStipendAllocation)
	r.Get("/api/stipend/allocations/{id}", handlers.GetStipendAllocationById)
	r.Put("/api/stipend/allocations/{id}", handlers.UpdateStipendAllocation)

	// ==================== Stipend History Endpoints ====================
	r.Get("/api/stipend/history", handlers.GetStipendHistory)
	r.Post("/api/stipend/history", handlers.CreateStipendHistory)
	r.Get("/api/stipend/history/student/{studentId}", handlers.GetStudentStipendHistory)

	// ==================== Report Endpoints ====================
	r.Get("/api/reports/students/summary", handlers.GenerateStudentSummary)
	r.Get("/api/reports/stipend/statistics", handlers.GenerateStipendStatistics)
	r.Get("/api/reports/students/by-college", handlers.GetStudentsByCollegeReport)
	r.Get("/api/reports/students/by-program", handlers.GetStudentsByProgramReport)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	log.Printf("Student Management Service starting on :%s", port)
	log.Println("Available endpoints:")
	log.Println("  - Students: /api/students")
	log.Println("  - Programs: /api/programs")
	log.Println("  - Colleges: /api/colleges")
	log.Println("  - Stipend: /api/stipend/*")
	log.Println("  - Reports: /api/reports/*")
	
	http.ListenAndServe(":"+port, r)
}
