package grpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "student_management_service/pb/student"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// StudentServiceClient wraps the gRPC student service client
type StudentServiceClient struct {
	conn   *grpc.ClientConn
	client pb.StudentServiceClient
}

// NewStudentServiceClient creates a new student service client
func NewStudentServiceClient() (*StudentServiceClient, error) {
	studentURL := os.Getenv("STUDENT_GRPC_URL")
	if studentURL == "" {
		studentURL = "localhost:50054"
	}

	log.Printf("Connecting to Student Service at %s", studentURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		studentURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to student service: %w", err)
	}

	client := pb.NewStudentServiceClient(conn)
	return &StudentServiceClient{
		conn:   conn,
		client: client,
	}, nil
}

// GetStudent retrieves a student by ID
func (sc *StudentServiceClient) GetStudent(ctx context.Context, studentID uint32) (*pb.StudentResponse, error) {
	req := &pb.GetStudentRequest{
		Id: uint32(studentID),
	}

	return sc.client.GetStudent(ctx, req)
}

// GetStudentByStudentId retrieves a student by their Student ID (RUB ID)
func (sc *StudentServiceClient) GetStudentByStudentId(ctx context.Context, studentID string) (*pb.StudentResponse, error) {
	req := &pb.GetStudentByStudentIdRequest{
		StudentId: studentID,
	}

	return sc.client.GetStudentByStudentId(ctx, req)
}

// ListStudents retrieves all students with optional status filter
func (sc *StudentServiceClient) ListStudents(ctx context.Context, status string) (*pb.ListStudentsResponse, error) {
	req := &pb.ListStudentsRequest{
		Status: status,
	}

	return sc.client.ListStudents(ctx, req)
}

// SearchStudents searches students by query
func (sc *StudentServiceClient) SearchStudents(ctx context.Context, query string) (*pb.ListStudentsResponse, error) {
	req := &pb.SearchStudentsRequest{
		Query: query,
	}

	return sc.client.SearchStudents(ctx, req)
}

// GetStudentsByProgram retrieves students by program ID
func (sc *StudentServiceClient) GetStudentsByProgram(ctx context.Context, programID uint32) (*pb.ListStudentsResponse, error) {
	req := &pb.GetByProgramRequest{
		ProgramId: programID,
	}

	return sc.client.GetStudentsByProgram(ctx, req)
}

// GetStudentsByCollege retrieves students by college ID
func (sc *StudentServiceClient) GetStudentsByCollege(ctx context.Context, collegeID uint32) (*pb.ListStudentsResponse, error) {
	req := &pb.GetByCollegeRequest{
		CollegeId: collegeID,
	}

	return sc.client.GetStudentsByCollege(ctx, req)
}

// Close closes the connection to the student service
func (sc *StudentServiceClient) Close() error {
	if sc.conn != nil {
		return sc.conn.Close()
	}
	return nil
}
