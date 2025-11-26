package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "student_management_service/pb/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserGRPCClient handles gRPC communication with User Service
type UserGRPCClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

// NewUserGRPCClient creates a new User Service gRPC client
func NewUserGRPCClient() (*UserGRPCClient, error) {
	// Get User Service gRPC address from environment
	address := os.Getenv("USER_GRPC_URL")
	if address == "" {
		address = "localhost:50052" // Default for local development
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
		return nil, fmt.Errorf("failed to connect to user service at %s: %w", address, err)
	}

	log.Printf("Connected to User Service gRPC at %s", address)

	return &UserGRPCClient{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}, nil
}

// GetUser retrieves user information by ID
func (c *UserGRPCClient) GetUser(ctx context.Context, userID uint32) (*pb.UserResponse, error) {
	req := &pb.GetUserRequest{
		Id: userID,
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp, nil
}

// ValidateToken validates a user token and returns user information
func (c *UserGRPCClient) ValidateToken(ctx context.Context, token string) (*pb.ValidateTokenResponse, error) {
	req := &pb.ValidateTokenRequest{
		Token: token,
	}

	resp, err := c.client.ValidateToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return resp, nil
}

// GetUserRole retrieves role information for a user
func (c *UserGRPCClient) GetUserRole(ctx context.Context, userID uint32) (*pb.UserRoleResponse, error) {
	req := &pb.GetUserRoleRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUserRole(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	return resp, nil
}

// Close closes the gRPC connection
func (c *UserGRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
