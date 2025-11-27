package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "student_management_service/pb/finance"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// FinanceClient wraps the gRPC finance service client
type FinanceClient struct {
	conn   *grpc.ClientConn
	client pb.StipendServiceClient
}

// NewFinanceClient creates a new finance service client
func NewFinanceClient() (*FinanceClient, error) {
	financeURL := os.Getenv("FINANCE_GRPC_URL")
	if financeURL == "" {
		financeURL = "localhost:50055"
	}

	log.Printf("Connecting to Finance Service at %s", financeURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		financeURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to finance service: %w", err)
	}

	client := pb.NewStipendServiceClient(conn)
	return &FinanceClient{
		conn:   conn,
		client: client,
	}, nil
}

// CalculateStipendWithDeductions calls the finance service to calculate stipend with deductions
func (fc *FinanceClient) CalculateStipendWithDeductions(
	ctx context.Context,
	studentID string,
	stipendType string,
	amount float64,
) (*pb.CalculationResponse, error) {
	req := &pb.CalculateStipendRequest{
		StudentId:   studentID,
		StipendType: stipendType,
		Amount:      amount,
	}

	return fc.client.CalculateStipendWithDeductions(ctx, req)
}

// CalculateMonthlyStipend calls the finance service to calculate monthly stipend
func (fc *FinanceClient) CalculateMonthlyStipend(
	ctx context.Context,
	studentID string,
	stipendType string,
	annualAmount float64,
) (*pb.CalculationResponse, error) {
	req := &pb.CalculateMonthlyStipendRequest{
		StudentId:    studentID,
		StipendType:  stipendType,
		AnnualAmount: annualAmount,
	}

	return fc.client.CalculateMonthlyStipend(ctx, req)
}

// CalculateAnnualStipend calls the finance service to calculate annual stipend
func (fc *FinanceClient) CalculateAnnualStipend(
	ctx context.Context,
	studentID string,
	stipendType string,
	amount float64,
) (*pb.CalculationResponse, error) {
	req := &pb.CalculateAnnualStipendRequest{
		StudentId:   studentID,
		StipendType: stipendType,
		Amount:      amount,
	}

	return fc.client.CalculateAnnualStipend(ctx, req)
}

// GetStudentStipends retrieves all stipends for a student
func (fc *FinanceClient) GetStudentStipends(
	ctx context.Context,
	studentID string,
) (*pb.StudentStipendsResponse, error) {
	req := &pb.GetStudentStipendsRequest{
		StudentId: studentID,
	}

	return fc.client.GetStudentStipends(ctx, req)
}

// CreateStipend creates a new stipend record
func (fc *FinanceClient) CreateStipend(
	ctx context.Context,
	studentID string,
	stipendType string,
	amount float64,
	paymentMethod string,
	journalNumber string,
	notes string,
) (*pb.StipendResponse, error) {
	req := &pb.CreateStipendRequest{
		StudentId:     studentID,
		StipendType:   stipendType,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		JournalNumber: journalNumber,
		Notes:         notes,
	}

	return fc.client.CreateStipend(ctx, req)
}

// GetStipend retrieves a specific stipend by ID
func (fc *FinanceClient) GetStipend(ctx context.Context, stipendID string) (*pb.StipendResponse, error) {
	req := &pb.GetStipendRequest{
		StipendId: stipendID,
	}

	return fc.client.GetStipend(ctx, req)
}

// UpdateStipendPaymentStatus updates the payment status of a stipend
func (fc *FinanceClient) UpdateStipendPaymentStatus(
	ctx context.Context,
	stipendID string,
	paymentStatus string,
) (*pb.StipendResponse, error) {
	req := &pb.UpdateStipendPaymentStatusRequest{
		StipendId:     stipendID,
		PaymentStatus: paymentStatus,
	}

	return fc.client.UpdateStipendPaymentStatus(ctx, req)
}

// Close closes the connection to the finance service
func (fc *FinanceClient) Close() error {
	if fc.conn != nil {
		return fc.conn.Close()
	}
	return nil
}
