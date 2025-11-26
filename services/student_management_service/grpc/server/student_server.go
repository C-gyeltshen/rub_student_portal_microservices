package server

import (
	"context"
	"fmt"
	"student_management_service/database"
	"student_management_service/models"
	pb "student_management_service/pb/student"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// StudentServer implements the StudentService gRPC server
type StudentServer struct {
	pb.UnimplementedStudentServiceServer
}

// GetStudent retrieves a student by ID
func (s *StudentServer) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.StudentResponse, error) {
	var student models.Student
	
	if err := database.DB.Preload("Program").Preload("College").First(&student, req.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("student with ID %d not found", req.Id)
		}
		return nil, fmt.Errorf("failed to fetch student: %w", err)
	}

	return convertStudentToProto(&student), nil
}

// GetStudentByStudentId retrieves a student by their Student ID (RUB ID)
func (s *StudentServer) GetStudentByStudentId(ctx context.Context, req *pb.GetStudentByStudentIdRequest) (*pb.StudentResponse, error) {
	var student models.Student
	
	if err := database.DB.Preload("Program").Preload("College").Where("student_id = ?", req.StudentId).First(&student).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("student with Student ID %s not found", req.StudentId)
		}
		return nil, fmt.Errorf("failed to fetch student: %w", err)
	}

	return convertStudentToProto(&student), nil
}

// CreateStudent creates a new student record
func (s *StudentServer) CreateStudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.StudentResponse, error) {
	student := models.Student{
		FirstName:           req.FirstName,
		LastName:            req.LastName,
		StudentID:           req.StudentId,
		CID:                 req.Cid,
		Email:               req.Email,
		PhoneNumber:         req.PhoneNumber,
		DateOfBirth:         req.DateOfBirth,
		Gender:              req.Gender,
		ProgramID:           uint(req.ProgramId),
		CollegeID:           uint(req.CollegeId),
		UserID:              uint(req.UserId),
		PermanentAddress:    req.PermanentAddress,
		CurrentAddress:      req.CurrentAddress,
		GuardianName:        req.GuardianName,
		GuardianPhoneNumber: req.GuardianPhone,
		EnrollmentDate:      req.AdmissionDate,
		Status:              req.EnrollmentStatus,
	}

	if err := database.DB.Create(&student).Error; err != nil {
		return nil, fmt.Errorf("failed to create student: %w", err)
	}

	// Reload with associations
	database.DB.Preload("Program").Preload("College").First(&student, student.ID)

	return convertStudentToProto(&student), nil
}

// UpdateStudent updates an existing student record
func (s *StudentServer) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.StudentResponse, error) {
	var student models.Student
	
	if err := database.DB.First(&student, req.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("student with ID %d not found", req.Id)
		}
		return nil, fmt.Errorf("failed to fetch student: %w", err)
	}

	// Update fields
	updates := map[string]interface{}{
		"first_name":            req.FirstName,
		"last_name":             req.LastName,
		"email":                 req.Email,
		"phone_number":          req.PhoneNumber,
		"gender":                req.Gender,
		"program_id":            req.ProgramId,
		"college_id":            req.CollegeId,
		"permanent_address":     req.PermanentAddress,
		"current_address":       req.CurrentAddress,
		"guardian_name":         req.GuardianName,
		"guardian_phone_number": req.GuardianPhone,
		"status":                req.EnrollmentStatus,
		"gpa":                   req.Gpa,
		"academic_standing":     req.AcademicStanding,
	}

	if err := database.DB.Model(&student).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update student: %w", err)
	}

	// Reload with associations
	database.DB.Preload("Program").Preload("College").First(&student, student.ID)

	return convertStudentToProto(&student), nil
}

// DeleteStudent soft deletes a student record
func (s *StudentServer) DeleteStudent(ctx context.Context, req *pb.DeleteStudentRequest) (*pb.DeleteResponse, error) {
	var student models.Student
	
	if err := database.DB.First(&student, req.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.DeleteResponse{
				Success: false,
				Message: fmt.Sprintf("student with ID %d not found", req.Id),
			}, nil
		}
		return nil, fmt.Errorf("failed to fetch student: %w", err)
	}

	if err := database.DB.Delete(&student).Error; err != nil {
		return &pb.DeleteResponse{
			Success: false,
			Message: fmt.Sprintf("failed to delete student: %v", err),
		}, nil
	}

	return &pb.DeleteResponse{
		Success: true,
		Message: "Student deleted successfully",
	}, nil
}

// ListStudents retrieves all students with optional status filter
func (s *StudentServer) ListStudents(ctx context.Context, req *pb.ListStudentsRequest) (*pb.ListStudentsResponse, error) {
	var students []models.Student
	query := database.DB.Preload("Program").Preload("College")

	if req.Status != "" {
		query = query.Where("enrollment_status = ?", req.Status)
	}

	if err := query.Find(&students).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch students: %w", err)
	}

	studentResponses := make([]*pb.StudentResponse, len(students))
	for i, student := range students {
		studentResponses[i] = convertStudentToProto(&student)
	}

	return &pb.ListStudentsResponse{
		Students: studentResponses,
		Total:    int32(len(students)),
	}, nil
}

// SearchStudents searches for students by query string
func (s *StudentServer) SearchStudents(ctx context.Context, req *pb.SearchStudentsRequest) (*pb.ListStudentsResponse, error) {
	var students []models.Student
	query := database.DB.Preload("Program").Preload("College")

	searchPattern := "%" + req.Query + "%"
	query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR student_id ILIKE ?",
		searchPattern, searchPattern, searchPattern, searchPattern)

	if err := query.Find(&students).Error; err != nil {
		return nil, fmt.Errorf("failed to search students: %w", err)
	}

	studentResponses := make([]*pb.StudentResponse, len(students))
	for i, student := range students {
		studentResponses[i] = convertStudentToProto(&student)
	}

	return &pb.ListStudentsResponse{
		Students: studentResponses,
		Total:    int32(len(students)),
	}, nil
}

// GetStudentsByProgram retrieves students by program ID
func (s *StudentServer) GetStudentsByProgram(ctx context.Context, req *pb.GetByProgramRequest) (*pb.ListStudentsResponse, error) {
	var students []models.Student
	
	if err := database.DB.Preload("Program").Preload("College").Where("program_id = ?", req.ProgramId).Find(&students).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch students by program: %w", err)
	}

	studentResponses := make([]*pb.StudentResponse, len(students))
	for i, student := range students {
		studentResponses[i] = convertStudentToProto(&student)
	}

	return &pb.ListStudentsResponse{
		Students: studentResponses,
		Total:    int32(len(students)),
	}, nil
}

// GetStudentsByCollege retrieves students by college ID
func (s *StudentServer) GetStudentsByCollege(ctx context.Context, req *pb.GetByCollegeRequest) (*pb.ListStudentsResponse, error) {
	var students []models.Student
	
	if err := database.DB.Preload("Program").Preload("College").Where("college_id = ?", req.CollegeId).Find(&students).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch students by college: %w", err)
	}

	studentResponses := make([]*pb.StudentResponse, len(students))
	for i, student := range students {
		studentResponses[i] = convertStudentToProto(&student)
	}

	return &pb.ListStudentsResponse{
		Students: studentResponses,
		Total:    int32(len(students)),
	}, nil
}

// CheckStipendEligibility checks if a student is eligible for stipend
func (s *StudentServer) CheckStipendEligibility(ctx context.Context, req *pb.StipendEligibilityRequest) (*pb.StipendEligibilityResponse, error) {
	var student models.Student
	
	if err := database.DB.Preload("Program").First(&student, req.StudentId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.StipendEligibilityResponse{
				Eligible: false,
				Reason:   "Student not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to fetch student: %w", err)
	}

	// Check eligibility criteria
	if student.Status != "active" {
		return &pb.StipendEligibilityResponse{
			Eligible: false,
			Reason:   "Student is not active",
		}, nil
	}

	// Preload college to check stipend policy
	database.DB.Preload("College").First(&student, req.StudentId)

	// Check financing type and college policy
	if student.FinancingType == "self-financed" {
		// Check if college allows self-financed students to get stipend
		if student.College.ID != 0 && !student.College.AllowSelfFinancedStipend {
			return &pb.StipendEligibilityResponse{
				Eligible: false,
				Reason:   "College does not allow self-financed students to receive stipend",
			}, nil
		}
	}
	// Scholarship students are always eligible for stipend (if other criteria met)

	if !student.Program.HasStipend {
		return &pb.StipendEligibilityResponse{
			Eligible: false,
			Reason:   "Program does not offer stipend",
		}, nil
	}

	return &pb.StipendEligibilityResponse{
		Eligible:     true,
		Reason:       "Student meets all eligibility criteria",
		Amount:       student.Program.StipendAmount,
		StipendType:  student.Program.StipendType,
	}, nil
}

// Helper function to convert models.Student to pb.StudentResponse
func convertStudentToProto(student *models.Student) *pb.StudentResponse {
	return &pb.StudentResponse{
		Id:               uint32(student.ID),
		FirstName:        student.FirstName,
		LastName:         student.LastName,
		StudentId:        student.StudentID,
		Cid:              student.CID,
		Email:            student.Email,
		PhoneNumber:      student.PhoneNumber,
		DateOfBirth:      student.DateOfBirth,
		Gender:           student.Gender,
		ProgramId:        uint32(student.ProgramID),
		CollegeId:        uint32(student.CollegeID),
		UserId:           uint32(student.UserID),
		PermanentAddress: student.PermanentAddress,
		CurrentAddress:   student.CurrentAddress,
		GuardianName:     student.GuardianName,
		GuardianPhone:    student.GuardianPhoneNumber,
		AdmissionDate:    student.EnrollmentDate,
		EnrollmentStatus: student.Status,
		AcademicStanding: student.AcademicStanding,
		CreatedAt:        timestamppb.New(student.CreatedAt),
		UpdatedAt:        timestamppb.New(student.UpdatedAt),
	}
}
