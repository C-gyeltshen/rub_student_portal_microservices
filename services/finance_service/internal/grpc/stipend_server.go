package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	pb "finance_service/pkg/pb"
	"finance_service/services"
)

// StipendServiceServer implements the StipendService gRPC server
type StipendServiceServer struct {
	pb.UnimplementedStipendServiceServer
	stipendService *services.StipendService
}

// NewStipendServiceServer creates a new StipendService gRPC server
func NewStipendServiceServer() *StipendServiceServer {
	return &StipendServiceServer{
		stipendService: services.NewStipendService(),
	}
}

// CalculateStipendWithDeductions calculates stipend with deductions
func (s *StipendServiceServer) CalculateStipendWithDeductions(ctx context.Context, req *pb.CalculateStipendRequest) (*pb.CalculationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	studentID, err := uuid.Parse(req.StudentId)
	if err != nil {
		return nil, fmt.Errorf("invalid student_id: %w", err)
	}

	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	result, err := s.stipendService.CalculateStipendWithDeductions(studentID, req.StipendType, req.Amount)
	if err != nil {
		return nil, fmt.Errorf("calculation failed: %w", err)
	}

	return s.convertCalculationResultToProto(result), nil
}

// CalculateMonthlyStipend calculates monthly stipend with deductions
func (s *StipendServiceServer) CalculateMonthlyStipend(ctx context.Context, req *pb.CalculateMonthlyStipendRequest) (*pb.CalculationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	studentID, err := uuid.Parse(req.StudentId)
	if err != nil {
		return nil, fmt.Errorf("invalid student_id: %w", err)
	}

	if req.AnnualAmount <= 0 {
		return nil, fmt.Errorf("annual_amount must be positive")
	}

	// Calculate monthly amount
	monthlyAmount := req.AnnualAmount / 12

	result, err := s.stipendService.CalculateStipendWithDeductions(studentID, req.StipendType, monthlyAmount)
	if err != nil {
		return nil, fmt.Errorf("monthly calculation failed: %w", err)
	}

	return s.convertCalculationResultToProto(result), nil
}

// CalculateAnnualStipend calculates annual stipend with deductions
func (s *StipendServiceServer) CalculateAnnualStipend(ctx context.Context, req *pb.CalculateAnnualStipendRequest) (*pb.CalculationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	studentID, err := uuid.Parse(req.StudentId)
	if err != nil {
		return nil, fmt.Errorf("invalid student_id: %w", err)
	}

	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	result, err := s.stipendService.CalculateStipendWithDeductions(studentID, req.StipendType, req.Amount)
	if err != nil {
		return nil, fmt.Errorf("annual calculation failed: %w", err)
	}

	return s.convertCalculationResultToProto(result), nil
}

// CreateStipend creates a new stipend record
func (s *StipendServiceServer) CreateStipend(ctx context.Context, req *pb.CreateStipendRequest) (*pb.StipendResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	studentID, err := uuid.Parse(req.StudentId)
	if err != nil {
		return nil, fmt.Errorf("invalid student_id: %w", err)
	}

	stipend, err := s.stipendService.CreateStipendForStudent(
		studentID,
		req.StipendType,
		req.Amount,
		req.PaymentMethod,
		req.JournalNumber,
		req.Notes,
	)
	if err != nil {
		log.Printf("Error creating stipend: %v", err)
		return nil, fmt.Errorf("failed to create stipend: %w", err)
	}

	// Convert models.Stipend to services.Stipend
	serviceStipend := &services.Stipend{
		ID:            stipend.ID,
		StudentID:     stipend.StudentID,
		Amount:        stipend.Amount,
		StipendType:   stipend.StipendType,
		PaymentDate:   stipend.PaymentDate,
		PaymentStatus: stipend.PaymentStatus,
		PaymentMethod: stipend.PaymentMethod,
		JournalNumber: stipend.JournalNumber,
		Notes:         stipend.Notes,
		CreatedAt:     stipend.CreatedAt,
		ModifiedAt:    stipend.ModifiedAt,
	}

	return s.convertStipendToProto(serviceStipend), nil
}

// GetStipend retrieves a stipend by ID
func (s *StipendServiceServer) GetStipend(ctx context.Context, req *pb.GetStipendRequest) (*pb.StipendResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	stipendID, err := uuid.Parse(req.StipendId)
	if err != nil {
		return nil, fmt.Errorf("invalid stipend_id: %w", err)
	}

	stipend, err := s.stipendService.GetStipendByID(stipendID)
	if err != nil {
		log.Printf("Error getting stipend: %v", err)
		return nil, fmt.Errorf("failed to get stipend: %w", err)
	}

	// Convert models.Stipend to services.Stipend
	serviceStipend := &services.Stipend{
		ID:            stipend.ID,
		StudentID:     stipend.StudentID,
		Amount:        stipend.Amount,
		StipendType:   stipend.StipendType,
		PaymentDate:   stipend.PaymentDate,
		PaymentStatus: stipend.PaymentStatus,
		PaymentMethod: stipend.PaymentMethod,
		JournalNumber: stipend.JournalNumber,
		Notes:         stipend.Notes,
		CreatedAt:     stipend.CreatedAt,
		ModifiedAt:    stipend.ModifiedAt,
	}

	return s.convertStipendToProto(serviceStipend), nil
}

// GetStudentStipends retrieves all stipends for a student
func (s *StipendServiceServer) GetStudentStipends(ctx context.Context, req *pb.GetStudentStipendsRequest) (*pb.StudentStipendsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	studentID, err := uuid.Parse(req.StudentId)
	if err != nil {
		return nil, fmt.Errorf("invalid student_id: %w", err)
	}

	limit := int(req.Limit)
	offset := int(req.Offset)

	stipends, total, err := s.stipendService.GetStudentStipendsWithPagination(studentID, limit, offset)
	if err != nil {
		log.Printf("Error getting student stipends: %v", err)
		return nil, fmt.Errorf("failed to get student stipends: %w", err)
	}

	protoStipends := make([]*pb.StipendResponse, len(stipends))
	for i, stipend := range stipends {
		protoStipends[i] = s.convertStipendToProto(stipend)
	}

	return &pb.StudentStipendsResponse{
		Stipends: protoStipends,
		Total:    int32(total),
	}, nil
}

// UpdateStipendPaymentStatus updates the payment status of a stipend
func (s *StipendServiceServer) UpdateStipendPaymentStatus(ctx context.Context, req *pb.UpdateStipendPaymentStatusRequest) (*pb.StipendResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	stipendID, err := uuid.Parse(req.StipendId)
	if err != nil {
		return nil, fmt.Errorf("invalid stipend_id: %w", err)
	}

	var paymentDate *time.Time
	if req.PaymentDate != "" {
		t, err := time.Parse(time.RFC3339, req.PaymentDate)
		if err != nil {
			return nil, fmt.Errorf("invalid payment_date format: %w", err)
		}
		paymentDate = &t
	}

	stipend, err := s.stipendService.UpdateStipendPaymentStatusWithReturn(stipendID, req.PaymentStatus, paymentDate)
	if err != nil {
		log.Printf("Error updating stipend payment status: %v", err)
		return nil, fmt.Errorf("failed to update stipend payment status: %w", err)
	}

	return s.convertStipendToProto(stipend), nil
}

// Helper functions

func (s *StipendServiceServer) convertStipendToProto(stipend *services.Stipend) *pb.StipendResponse {
	var paymentDate int64
	if stipend.PaymentDate != nil {
		paymentDate = stipend.PaymentDate.Unix()
	}

	return &pb.StipendResponse{
		Id:            stipend.ID.String(),
		StudentId:     stipend.StudentID.String(),
		Amount:        stipend.Amount,
		StipendType:   stipend.StipendType,
		PaymentDate:   paymentDate,
		PaymentStatus: stipend.PaymentStatus,
		PaymentMethod: stipend.PaymentMethod,
		JournalNumber: stipend.JournalNumber,
		Notes:         stipend.Notes,
		CreatedAt:     stipend.CreatedAt.Unix(),
		ModifiedAt:    stipend.ModifiedAt.Unix(),
	}
}

func (s *StipendServiceServer) convertCalculationResultToProto(result *services.StipendCalculationResult) *pb.CalculationResponse {
	deductions := make([]*pb.DeductionDetail, len(result.Deductions))
	for i, d := range result.Deductions {
		deductions[i] = &pb.DeductionDetail{
			RuleId:        d.RuleID.String(),
			RuleName:      d.RuleName,
			DeductionType: d.DeductionType,
			Amount:        d.Amount,
			Description:   d.Description,
			IsOptional:    d.IsOptional,
		}
	}

	return &pb.CalculationResponse{
		BaseStipendAmount: result.BaseStipendAmount,
		TotalDeductions:   result.TotalDeductions,
		NetStipendAmount:  result.NetStipendAmount,
		Deductions:        deductions,
	}
}
