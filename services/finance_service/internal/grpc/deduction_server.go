package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	pb "finance_service/pkg/pb"
	"finance_service/services"
)

// DeductionServiceServer implements the DeductionService gRPC server
type DeductionServiceServer struct {
	pb.UnimplementedDeductionServiceServer
	deductionService *services.DeductionService
}

// NewDeductionServiceServer creates a new DeductionService gRPC server
func NewDeductionServiceServer() *DeductionServiceServer {
	return &DeductionServiceServer{
		deductionService: services.NewDeductionService(),
	}
}

// CreateDeductionRule creates a new deduction rule
func (s *DeductionServiceServer) CreateDeductionRule(ctx context.Context, req *pb.CreateDeductionRuleRequest) (*pb.DeductionRuleResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	rule, err := s.deductionService.CreateDeductionRule(
		req.RuleName,
		req.DeductionType,
		req.DefaultAmount,
		req.MinAmount,
		req.MaxAmount,
		req.Frequency,
		req.IsMandatory,
		req.ApplicableTo,
		int(req.Priority),
		req.Description,
	)
	if err != nil {
		log.Printf("Error creating deduction rule: %v", err)
		return nil, fmt.Errorf("failed to create deduction rule: %w", err)
	}

	return s.convertRuleToProto(rule), nil
}

// GetDeductionRule retrieves a deduction rule by ID
func (s *DeductionServiceServer) GetDeductionRule(ctx context.Context, req *pb.GetDeductionRuleRequest) (*pb.DeductionRuleResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	ruleID, err := uuid.Parse(req.RuleId)
	if err != nil {
		return nil, fmt.Errorf("invalid rule_id: %w", err)
	}

	rule, err := s.deductionService.GetDeductionRuleByID(ruleID)
	if err != nil {
		log.Printf("Error getting deduction rule: %v", err)
		return nil, fmt.Errorf("failed to get deduction rule: %w", err)
	}

	return s.convertRuleToProto(rule), nil
}

// ListDeductionRules lists all deduction rules with optional filters
func (s *DeductionServiceServer) ListDeductionRules(ctx context.Context, req *pb.ListDeductionRulesRequest) (*pb.DeductionRulesResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	limit := int(req.Limit)
	offset := int(req.Offset)

	rules, total, err := s.deductionService.ListDeductionRulesWithPagination(req.ApplicableTo, req.IsActive, limit, offset)
	if err != nil {
		log.Printf("Error listing deduction rules: %v", err)
		return nil, fmt.Errorf("failed to list deduction rules: %w", err)
	}

	protoRules := make([]*pb.DeductionRuleResponse, len(rules))
	for i, rule := range rules {
		protoRules[i] = s.convertRuleToProto(rule)
	}

	return &pb.DeductionRulesResponse{
		Rules: protoRules,
		Total: int32(total),
	}, nil
}

// ApplyDeductions applies deductions to a stipend
func (s *DeductionServiceServer) ApplyDeductions(ctx context.Context, req *pb.ApplyDeductionsRequest) (*pb.DeductionsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	studentID, err := uuid.Parse(req.StudentId)
	if err != nil {
		return nil, fmt.Errorf("invalid student_id: %w", err)
	}

	stipendID, err := uuid.Parse(req.StipendId)
	if err != nil {
		return nil, fmt.Errorf("invalid stipend_id: %w", err)
	}

	ruleIDs := make([]uuid.UUID, len(req.DeductionRuleIds))
	for i, id := range req.DeductionRuleIds {
		ruleID, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("invalid rule_id at index %d: %w", i, err)
		}
		ruleIDs[i] = ruleID
	}

	deductions, totalAmount, err := s.deductionService.ApplyDeductions(
		studentID,
		stipendID,
		req.StipendType,
		req.BaseAmount,
		ruleIDs,
	)
	if err != nil {
		log.Printf("Error applying deductions: %v", err)
		return nil, fmt.Errorf("failed to apply deductions: %w", err)
	}

	protoDeductions := make([]*pb.DeductionResponse, len(deductions))
	for i, deduction := range deductions {
		protoDeductions[i] = s.convertDeductionToProto(deduction)
	}

	return &pb.DeductionsResponse{
		Deductions:  protoDeductions,
		TotalAmount: totalAmount,
		Total:       int32(len(deductions)),
	}, nil
}

// CreateDeduction creates a new deduction record
func (s *DeductionServiceServer) CreateDeduction(ctx context.Context, req *pb.CreateDeductionRequest) (*pb.DeductionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	studentID, err := uuid.Parse(req.StudentId)
	if err != nil {
		return nil, fmt.Errorf("invalid student_id: %w", err)
	}

	ruleID, err := uuid.Parse(req.DeductionRuleId)
	if err != nil {
		return nil, fmt.Errorf("invalid deduction_rule_id: %w", err)
	}

	stipendID, err := uuid.Parse(req.StipendId)
	if err != nil {
		return nil, fmt.Errorf("invalid stipend_id: %w", err)
	}

	deduction, err := s.deductionService.CreateDeduction(
		studentID,
		ruleID,
		stipendID,
		req.Amount,
		req.DeductionType,
		req.Description,
	)
	if err != nil {
		log.Printf("Error creating deduction: %v", err)
		return nil, fmt.Errorf("failed to create deduction: %w", err)
	}

	return s.convertDeductionToProto(deduction), nil
}

// GetDeduction retrieves a deduction by ID
func (s *DeductionServiceServer) GetDeduction(ctx context.Context, req *pb.GetDeductionRequest) (*pb.DeductionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	deductionID, err := uuid.Parse(req.DeductionId)
	if err != nil {
		return nil, fmt.Errorf("invalid deduction_id: %w", err)
	}

	deduction, err := s.deductionService.GetDeductionByID(deductionID)
	if err != nil {
		log.Printf("Error getting deduction: %v", err)
		return nil, fmt.Errorf("failed to get deduction: %w", err)
	}

	return s.convertDeductionToProto(deduction), nil
}

// GetStipendDeductions retrieves all deductions for a stipend
func (s *DeductionServiceServer) GetStipendDeductions(ctx context.Context, req *pb.GetStipendDeductionsRequest) (*pb.DeductionsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	stipendID, err := uuid.Parse(req.StipendId)
	if err != nil {
		return nil, fmt.Errorf("invalid stipend_id: %w", err)
	}

	limit := int(req.Limit)
	offset := int(req.Offset)

	deductions, totalAmount, err := s.deductionService.GetStipendDeductionsWithPagination(stipendID, limit, offset)
	if err != nil {
		log.Printf("Error getting stipend deductions: %v", err)
		return nil, fmt.Errorf("failed to get stipend deductions: %w", err)
	}

	protoDeductions := make([]*pb.DeductionResponse, len(deductions))
	for i, deduction := range deductions {
		protoDeductions[i] = s.convertDeductionToProto(deduction)
	}

	return &pb.DeductionsResponse{
		Deductions:  protoDeductions,
		TotalAmount: totalAmount,
		Total:       int32(len(deductions)),
	}, nil
}

// GetStudentDeductions retrieves all deductions for a student
func (s *DeductionServiceServer) GetStudentDeductions(ctx context.Context, req *pb.GetStudentDeductionsRequest) (*pb.DeductionsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	studentID, err := uuid.Parse(req.StudentId)
	if err != nil {
		return nil, fmt.Errorf("invalid student_id: %w", err)
	}

	limit := int(req.Limit)
	offset := int(req.Offset)

	deductions, totalAmount, err := s.deductionService.GetStudentDeductionsWithPagination(studentID, limit, offset)
	if err != nil {
		log.Printf("Error getting student deductions: %v", err)
		return nil, fmt.Errorf("failed to get student deductions: %w", err)
	}

	protoDeductions := make([]*pb.DeductionResponse, len(deductions))
	for i, deduction := range deductions {
		protoDeductions[i] = s.convertDeductionToProto(deduction)
	}

	return &pb.DeductionsResponse{
		Deductions:  protoDeductions,
		TotalAmount: totalAmount,
		Total:       int32(len(deductions)),
	}, nil
}

// Helper functions

func (s *DeductionServiceServer) convertRuleToProto(rule *services.DeductionRule) *pb.DeductionRuleResponse {
	return &pb.DeductionRuleResponse{
		Id:                rule.ID.String(),
		RuleName:          rule.RuleName,
		DeductionType:     rule.DeductionType,
		DefaultAmount:     rule.DefaultAmount,
		MinAmount:         rule.MinAmount,
		MaxAmount:         rule.MaxAmount,
		IsMandatory:       !rule.IsOptional,
		Priority:          int32(rule.Priority),
		Description:       rule.Description,
		IsActive:          rule.IsActive,
		CreatedAt:         rule.CreatedAt.Unix(),
		ModifiedAt:        rule.ModifiedAt.Unix(),
	}
}

func (s *DeductionServiceServer) convertDeductionToProto(deduction *services.Deduction) *pb.DeductionResponse {
	var approvalDate int64
	if deduction.ApprovalDate != nil {
		approvalDate = deduction.ApprovalDate.Unix()
	}

	var approvedBy string
	if deduction.ApprovedBy != nil {
		approvedBy = deduction.ApprovedBy.String()
	}

	var transactionID string
	if deduction.TransactionID != nil {
		transactionID = deduction.TransactionID.String()
	}

	return &pb.DeductionResponse{
		Id:               deduction.ID.String(),
		StudentId:        deduction.StudentID.String(),
		DeductionRuleId:  deduction.DeductionRuleID.String(),
		StipendId:        deduction.StipendID.String(),
		Amount:           deduction.Amount,
		DeductionType:    deduction.DeductionType,
		Description:      deduction.Description,
		DeductionDate:    deduction.DeductionDate.Unix(),
		ProcessingStatus: deduction.ProcessingStatus,
		ApprovedBy:       approvedBy,
		ApprovalDate:     approvalDate,
		RejectionReason:  deduction.RejectionReason,
		TransactionId:    transactionID,
		CreatedAt:        deduction.CreatedAt.Unix(),
		ModifiedAt:       deduction.ModifiedAt.Unix(),
	}
}
