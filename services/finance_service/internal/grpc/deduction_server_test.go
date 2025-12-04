package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "finance_service/pkg/pb"
)

// This file contains integration tests for the Deduction Service gRPC APIs
// To run these tests, ensure the gRPC server is running on port 50051

func TestDeductionService_CreateDeductionRule(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewDeductionServiceClient(conn)

	req := &pb.CreateDeductionRuleRequest{
		RuleName:      "Test Hostel Fee " + GetUniqueTestName(""),
		DeductionType: "hostel_test",
		DefaultAmount: 5000.00,
		MinAmount:     4000.00,
		MaxAmount:     6000.00,
		Frequency:     "Monthly",
		IsMandatory:   true,
		ApplicableTo:  "All",
		Priority:      100,
		Description:   "Test hostel deduction rule",
	}

	ctx := context.Background()
	resp, err := client.CreateDeductionRule(ctx, req)
	if err != nil {
		t.Fatalf("CreateDeductionRule failed: %v", err)
	}

	if resp.RuleName != req.RuleName {
		t.Errorf("Rule name mismatch: expected %s, got %s", req.RuleName, resp.RuleName)
	}

	if resp.DefaultAmount != req.DefaultAmount {
		t.Errorf("Default amount mismatch: expected %.2f, got %.2f", req.DefaultAmount, resp.DefaultAmount)
	}

	if !resp.IsActive {
		t.Error("Created rule should be active")
	}

	t.Logf("Created deduction rule with ID: %s", resp.Id)
}

func TestDeductionService_GetDeductionRule(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewDeductionServiceClient(conn)

	// Create a rule first
	createReq := &pb.CreateDeductionRuleRequest{
		RuleName:      "Test Electricity " + GetUniqueTestName(""),
		DeductionType: "electricity_test",
		DefaultAmount: 1000.00,
		MinAmount:     500.00,
		MaxAmount:     1500.00,
		Frequency:     "Monthly",
		IsMandatory:   true,
		ApplicableTo:  "All",
		Priority:      80,
		Description:   "Test electricity deduction",
	}

	ctx := context.Background()
	createResp, err := client.CreateDeductionRule(ctx, createReq)
	if err != nil {
		t.Fatalf("CreateDeductionRule failed: %v", err)
	}

	// Now retrieve it
	getReq := &pb.GetDeductionRuleRequest{
		RuleId: createResp.Id,
	}

	getResp, err := client.GetDeductionRule(ctx, getReq)
	if err != nil {
		t.Fatalf("GetDeductionRule failed: %v", err)
	}

	if getResp.Id != createResp.Id {
		t.Errorf("Rule ID mismatch")
	}

	if getResp.RuleName != createReq.RuleName {
		t.Errorf("Rule name mismatch")
	}

	t.Logf("Retrieved deduction rule: %s", getResp.RuleName)
}

func TestDeductionService_ListDeductionRules(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewDeductionServiceClient(conn)

	req := &pb.ListDeductionRulesRequest{
		IsActive: true,
		Limit:    10,
		Offset:   0,
	}

	ctx := context.Background()
	resp, err := client.ListDeductionRules(ctx, req)
	if err != nil {
		t.Fatalf("ListDeductionRules failed: %v", err)
	}

	if resp.Total < 0 {
		t.Errorf("Total should not be negative: %d", resp.Total)
	}

	if int32(len(resp.Rules)) > resp.Total {
		t.Errorf("Number of rules should not exceed total")
	}

	t.Logf("Listed %d deduction rules (total: %d)", len(resp.Rules), resp.Total)
}

func TestDeductionService_CreateDeduction(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	deductionClient := pb.NewDeductionServiceClient(conn)
	stipendClient := pb.NewStipendServiceClient(conn)

	ctx := context.Background()

	// Create a stipend first
	studentID := GetTestStudentID()
	stipendReq := &pb.CreateStipendRequest{
		StudentId:     studentID,
		StipendType:   "full-scholarship",
		Amount:        100000.00,
		PaymentMethod: "Bank_transfer",
		JournalNumber: "JN-DED-" + GetUniqueTestName("001"),
	}

	stipendResp, err := stipendClient.CreateStipend(ctx, stipendReq)
	if err != nil {
		t.Fatalf("CreateStipend failed: %v", err)
	}

	// Create a deduction rule
	ruleReq := &pb.CreateDeductionRuleRequest{
		RuleName:      "Test Mess Fee " + GetUniqueTestName(""),
		DeductionType: "mess_fees_test",
		DefaultAmount: 3000.00,
		MinAmount:     2500.00,
		MaxAmount:     3500.00,
		Frequency:     "Monthly",
		IsMandatory:   true,
		ApplicableTo:  "All",
		Priority:      90,
		Description:   "Test mess fee",
	}

	ruleResp, err := deductionClient.CreateDeductionRule(ctx, ruleReq)
	if err != nil {
		t.Fatalf("CreateDeductionRule failed: %v", err)
	}

	// Create a deduction
	deductionReq := &pb.CreateDeductionRequest{
		StudentId:       studentID,
		DeductionRuleId: ruleResp.Id,
		StipendId:       stipendResp.Id,
		Amount:          3000.00,
		DeductionType:   "mess_fees_test",
		Description:     "Test mess fee deduction",
	}

	deductionResp, err := deductionClient.CreateDeduction(ctx, deductionReq)
	if err != nil {
		t.Fatalf("CreateDeduction failed: %v", err)
	}

	if deductionResp.Amount != 3000.00 {
		t.Errorf("Amount mismatch: expected 3000.00, got %.2f", deductionResp.Amount)
	}

	if deductionResp.ProcessingStatus != "Pending" {
		t.Errorf("Expected status Pending, got %s", deductionResp.ProcessingStatus)
	}

	t.Logf("Created deduction with ID: %s, Amount: %.2f", deductionResp.Id, deductionResp.Amount)
}

func TestDeductionService_GetDeduction(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	deductionClient := pb.NewDeductionServiceClient(conn)
	stipendClient := pb.NewStipendServiceClient(conn)

	ctx := context.Background()

	// Create a stipend and deduction
	studentID := GetTestStudentID()
	stipendReq := &pb.CreateStipendRequest{
		StudentId:     studentID,
		StipendType:   "self-funded",
		Amount:        120000.00,
		PaymentMethod: "E-payment",
		JournalNumber: "JN-DED-" + GetUniqueTestName("002"),
	}

	stipendResp, err := stipendClient.CreateStipend(ctx, stipendReq)
	if err != nil {
		t.Fatalf("CreateStipend failed: %v", err)
	}

	// Create a deduction rule
	ruleReq := &pb.CreateDeductionRuleRequest{
		RuleName:      "Test Water Bill " + GetUniqueTestName(""),
		DeductionType: "water_test",
		DefaultAmount: 500.00,
		MinAmount:     300.00,
		MaxAmount:     700.00,
		Frequency:     "Monthly",
		IsMandatory:   true,
		ApplicableTo:  "All",
		Priority:      70,
		Description:   "Test water bill",
	}

	ruleResp, err := deductionClient.CreateDeductionRule(ctx, ruleReq)
	if err != nil {
		t.Fatalf("CreateDeductionRule failed: %v", err)
	}

	// Create a deduction
	deductionReq := &pb.CreateDeductionRequest{
		StudentId:       studentID,
		DeductionRuleId: ruleResp.Id,
		StipendId:       stipendResp.Id,
		Amount:          500.00,
		DeductionType:   "water_test",
		Description:     "Test water bill deduction",
	}

	deductionResp, err := deductionClient.CreateDeduction(ctx, deductionReq)
	if err != nil {
		t.Fatalf("CreateDeduction failed: %v", err)
	}

	// Get the deduction
	getReq := &pb.GetDeductionRequest{
		DeductionId: deductionResp.Id,
	}

	getResp, err := deductionClient.GetDeduction(ctx, getReq)
	if err != nil {
		t.Fatalf("GetDeduction failed: %v", err)
	}

	if getResp.Id != deductionResp.Id {
		t.Errorf("Deduction ID mismatch")
	}

	if getResp.Amount != 500.00 {
		t.Errorf("Amount mismatch: expected 500.00, got %.2f", getResp.Amount)
	}

	t.Logf("Retrieved deduction: %s with amount %.2f", getResp.Id, getResp.Amount)
}

func TestDeductionService_GetStipendDeductions(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	deductionClient := pb.NewDeductionServiceClient(conn)
	stipendClient := pb.NewStipendServiceClient(conn)

	ctx := context.Background()

	// Create a stipend
	studentID := GetTestStudentID()
	stipendReq := &pb.CreateStipendRequest{
		StudentId:     studentID,
		StipendType:   "full-scholarship",
		Amount:        150000.00,
		PaymentMethod: "Bank_transfer",
		JournalNumber: "JN-DED-" + GetUniqueTestName("003"),
	}

	stipendResp, err := stipendClient.CreateStipend(ctx, stipendReq)
	if err != nil {
		t.Fatalf("CreateStipend failed: %v", err)
	}

	// List deductions for the stipend
	getReq := &pb.GetStipendDeductionsRequest{
		StipendId: stipendResp.Id,
		Limit:     10,
		Offset:    0,
	}

	resp, err := deductionClient.GetStipendDeductions(ctx, getReq)
	if err != nil {
		t.Fatalf("GetStipendDeductions failed: %v", err)
	}

	t.Logf("Retrieved %d deductions for stipend (total: %d)", len(resp.Deductions), resp.Total)
}

func TestDeductionService_ApplyDeductions(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Skipf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	deductionClient := pb.NewDeductionServiceClient(conn)
	stipendClient := pb.NewStipendServiceClient(conn)

	ctx := context.Background()

	// Create a stipend
	studentID := GetTestStudentID()
	stipendReq := &pb.CreateStipendRequest{
		StudentId:     studentID,
		StipendType:   "full-scholarship",
		Amount:        200000.00,
		PaymentMethod: "Bank_transfer",
		JournalNumber: "JN-DED-" + GetUniqueTestName("APPLY"),
	}

	stipendResp, err := stipendClient.CreateStipend(ctx, stipendReq)
	if err != nil {
		t.Fatalf("CreateStipend failed: %v", err)
	}

	// Apply deductions to the stipend
	applyReq := &pb.ApplyDeductionsRequest{
		StudentId:   studentID,
		StipendId:   stipendResp.Id,
		StipendType: "full-scholarship",
		BaseAmount:  200000.00,
	}

	resp, err := deductionClient.ApplyDeductions(ctx, applyReq)
	if err != nil {
		t.Fatalf("ApplyDeductions failed: %v", err)
	}

	if resp.TotalAmount < 0 {
		t.Errorf("Total deduction amount should not be negative: %.2f", resp.TotalAmount)
	}

	t.Logf("Applied %d deductions with total amount: %.2f", resp.Total, resp.TotalAmount)
}
