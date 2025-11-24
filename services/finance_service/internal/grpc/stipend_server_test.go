package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "finance_service/pkg/pb"
)

// This file contains integration tests for the gRPC APIs
// To run these tests, ensure the gRPC server is running on port 50051

func TestStipendService_CalculateStipendWithDeductions(t *testing.T) {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v (ensure server is running)", err)
	}
	defer conn.Close()

	client := pb.NewStipendServiceClient(conn)

	// Create a test request
	studentID := GetTestStudentID()
	req := &pb.CalculateStipendRequest{
		StudentId:     studentID,
		StipendType:   "full-scholarship",
		Amount:        100000.00,
		PaymentMethod: "Bank_transfer",
		JournalNumber: "JN-" + GetUniqueTestName("2024"),
		Notes:         "Test stipend calculation",
	}

	// Call the service
	ctx := context.Background()
	resp, err := client.CalculateStipendWithDeductions(ctx, req)
	if err != nil {
		t.Fatalf("CalculateStipendWithDeductions failed: %v", err)
	}

	// Assertions
	if resp.BaseStipendAmount != 100000.00 {
		t.Errorf("Expected base amount 100000.00, got %v", resp.BaseStipendAmount)
	}

	if resp.TotalDeductions < 0 {
		t.Errorf("Total deductions should not be negative: %v", resp.TotalDeductions)
	}

	if resp.NetStipendAmount < 0 {
		t.Errorf("Net stipend amount should not be negative: %v", resp.NetStipendAmount)
	}

	if resp.NetStipendAmount > resp.BaseStipendAmount {
		t.Errorf("Net stipend should not exceed base amount")
	}

	t.Logf("Calculation Result: Base=%.2f, Deductions=%.2f, Net=%.2f",
		resp.BaseStipendAmount, resp.TotalDeductions, resp.NetStipendAmount)
	t.Logf("Number of deductions applied: %d", len(resp.Deductions))
}

func TestStipendService_CalculateMonthlyStipend(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewStipendServiceClient(conn)

	studentID := GetTestStudentID()
	annualAmount := 1200000.00

	req := &pb.CalculateMonthlyStipendRequest{
		StudentId:     studentID,
		StipendType:   "self-funded",
		AnnualAmount:  annualAmount,
	}

	ctx := context.Background()
	resp, err := client.CalculateMonthlyStipend(ctx, req)
	if err != nil {
		t.Fatalf("CalculateMonthlyStipend failed: %v", err)
	}

	expectedMonthly := annualAmount / 12
	if resp.BaseStipendAmount != expectedMonthly {
		t.Errorf("Expected base amount %.2f, got %.2f", expectedMonthly, resp.BaseStipendAmount)
	}

	t.Logf("Monthly Calculation Result: Base=%.2f, Deductions=%.2f, Net=%.2f",
		resp.BaseStipendAmount, resp.TotalDeductions, resp.NetStipendAmount)
}

func TestStipendService_CalculateAnnualStipend(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewStipendServiceClient(conn)

	studentID := GetTestStudentID()
	annualAmount := 1200000.00

	req := &pb.CalculateAnnualStipendRequest{
		StudentId:   studentID,
		StipendType: "full-scholarship",
		Amount:      annualAmount,
	}

	ctx := context.Background()
	resp, err := client.CalculateAnnualStipend(ctx, req)
	if err != nil {
		t.Fatalf("CalculateAnnualStipend failed: %v", err)
	}

	if resp.BaseStipendAmount != annualAmount {
		t.Errorf("Expected base amount %.2f, got %.2f", annualAmount, resp.BaseStipendAmount)
	}

	t.Logf("Annual Calculation Result: Base=%.2f, Deductions=%.2f, Net=%.2f",
		resp.BaseStipendAmount, resp.TotalDeductions, resp.NetStipendAmount)
}

func TestStipendService_CreateStipend(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewStipendServiceClient(conn)

	studentID := GetTestStudentID()
	req := &pb.CreateStipendRequest{
		StudentId:     studentID,
		StipendType:   "full-scholarship",
		Amount:        50000.00,
		PaymentMethod: "Bank_transfer",
		JournalNumber: "JN-TEST-" + GetUniqueTestName("001"),
		Notes:         "Test stipend creation",
	}

	ctx := context.Background()
	resp, err := client.CreateStipend(ctx, req)
	if err != nil {
		t.Fatalf("CreateStipend failed: %v", err)
	}

	if resp.StudentId != studentID {
		t.Errorf("Student ID mismatch: expected %s, got %s", studentID, resp.StudentId)
	}

	if resp.Amount != 50000.00 {
		t.Errorf("Amount mismatch: expected 50000.00, got %v", resp.Amount)
	}

	if resp.PaymentStatus != "Pending" {
		t.Errorf("Expected payment status Pending, got %s", resp.PaymentStatus)
	}

	t.Logf("Created stipend with ID: %s", resp.Id)
}

func TestStipendService_GetStipend(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewStipendServiceClient(conn)

	// First create a stipend
	studentID := GetTestStudentID()
	createReq := &pb.CreateStipendRequest{
		StudentId:     studentID,
		StipendType:   "self-funded",
		Amount:        75000.00,
		PaymentMethod: "E-payment",
		JournalNumber: "JN-TEST-" + GetUniqueTestName("002"),
	}

	ctx := context.Background()
	createResp, err := client.CreateStipend(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateStipend failed: %v", err)
	}

	// Then retrieve it
	getReq := &pb.GetStipendRequest{
		StipendId: createResp.Id,
	}

	getResp, err := client.GetStipend(ctx, getReq)
	if err != nil {
		t.Fatalf("GetStipend failed: %v", err)
	}

	if getResp.Id != createResp.Id {
		t.Errorf("Stipend ID mismatch")
	}

	if getResp.Amount != 75000.00 {
		t.Errorf("Amount mismatch")
	}

	t.Logf("Retrieved stipend: %s", getResp.Id)
}

func TestStipendService_UpdateStipendPaymentStatus(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewStipendServiceClient(conn)

	// Create a stipend first
	studentID := GetTestStudentID()
	createReq := &pb.CreateStipendRequest{
		StudentId:     studentID,
		StipendType:   "full-scholarship",
		Amount:        60000.00,
		PaymentMethod: "Bank_transfer",
		JournalNumber: "JN-TEST-" + GetUniqueTestName("003"),
	}

	ctx := context.Background()
	createResp, err := client.CreateStipend(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateStipend failed: %v", err)
	}

	// Update payment status
	updateReq := &pb.UpdateStipendPaymentStatusRequest{
		StipendId:     createResp.Id,
		PaymentStatus: "Processed",
		PaymentDate:   "2024-11-25T10:30:00Z",
	}

	updateResp, err := client.UpdateStipendPaymentStatus(ctx, updateReq)
	if err != nil {
		t.Fatalf("UpdateStipendPaymentStatus failed: %v", err)
	}

	if updateResp.PaymentStatus != "Processed" {
		t.Errorf("Expected payment status Processed, got %s", updateResp.PaymentStatus)
	}

	t.Logf("Updated stipend payment status to: %s", updateResp.PaymentStatus)
}
