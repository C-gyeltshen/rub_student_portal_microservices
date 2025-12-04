package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "student_management_service/pb/banking"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// BankingGRPCClient handles gRPC communication with Banking Service
type BankingGRPCClient struct {
	conn   *grpc.ClientConn
	client pb.BankingServiceClient
}

// NewBankingGRPCClient creates a new Banking Service gRPC client
func NewBankingGRPCClient() (*BankingGRPCClient, error) {
	// Get Banking Service gRPC address from environment
	address := os.Getenv("BANKING_GRPC_URL")
	if address == "" {
		address = "localhost:50053" // Default for local development
	}

	// Create gRPC connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to banking service at %s: %w", address, err)
	}

	log.Printf("Connected to Banking Service gRPC at %s", address)

	return &BankingGRPCClient{
		conn:   conn,
		client: pb.NewBankingServiceClient(conn),
	}, nil
}

// GetStudentBankDetails retrieves bank details for a student
func (c *BankingGRPCClient) GetStudentBankDetails(ctx context.Context, studentID uint32) (*pb.BankDetailsResponse, error) {
	req := &pb.GetBankDetailsRequest{
		StudentId: studentID,
	}

	resp, err := c.client.GetStudentBankDetails(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get student bank details: %w", err)
	}

	return resp, nil
}

// UpsertBankDetails creates or updates bank details for a student
func (c *BankingGRPCClient) UpsertBankDetails(ctx context.Context, studentID, bankID uint32, accountNumber, accountHolderName string) (*pb.BankDetailsResponse, error) {
	req := &pb.UpsertBankDetailsRequest{
		StudentId:         studentID,
		BankId:            bankID,
		AccountNumber:     accountNumber,
		AccountHolderName: accountHolderName,
	}

	resp, err := c.client.UpsertBankDetails(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert bank details: %w", err)
	}

	return resp, nil
}

// VerifyBankAccount verifies if a bank account is valid
func (c *BankingGRPCClient) VerifyBankAccount(ctx context.Context, accountNumber string, bankID uint32) (*pb.VerifyBankAccountResponse, error) {
	req := &pb.VerifyBankAccountRequest{
		AccountNumber: accountNumber,
		BankId:        bankID,
	}

	resp, err := c.client.VerifyBankAccount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify bank account: %w", err)
	}

	return resp, nil
}

// Close closes the gRPC connection
func (c *BankingGRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
