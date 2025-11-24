# Finance Service gRPC APIs

## Overview

The Finance Service provides two main gRPC services for stipend calculation and deduction management:

1. **StipendService** - Calculate stipends with deductions and manage stipend records
2. **DeductionService** - Manage deduction rules and apply deductions to stipends

## Server Configuration

### gRPC Server

- **Port**: `50051` (default)
- **Environment Variable**: `GRPC_PORT`

### REST API Server

- **Port**: `8084` (default)
- **Environment Variable**: `PORT`

Both servers run concurrently when the Finance Service starts.

## Services

### 1. StipendService

The StipendService provides RPCs for calculating stipends and managing stipend records.

#### RPCs

##### CalculateStipendWithDeductions

Calculates a stipend amount after applying applicable deductions without creating a record.

```protobuf
rpc CalculateStipendWithDeductions(CalculateStipendRequest) returns (CalculationResponse);
```

**Request:**

```protobuf
message CalculateStipendRequest {
  string student_id = 1;
  string stipend_type = 2; // "full-scholarship" or "self-funded"
  double amount = 3;
  string payment_method = 4;
  string journal_number = 5;
  string notes = 6;
}
```

**Response:**

```protobuf
message CalculationResponse {
  double base_stipend_amount = 1;
  double total_deductions = 2;
  double net_stipend_amount = 3;
  repeated DeductionDetail deductions = 4;
}
```

**Example (Go):**

```go
client := pb.NewStipendServiceClient(conn)

req := &pb.CalculateStipendRequest{
    StudentId:     "550e8400-e29b-41d4-a716-446655440000",
    StipendType:   "full-scholarship",
    Amount:        100000.00,
    PaymentMethod: "Bank_transfer",
    JournalNumber: "JN-001-2024",
    Notes:         "January stipend",
}

resp, err := client.CalculateStipendWithDeductions(context.Background(), req)
```

##### CalculateMonthlyStipend

Calculates monthly stipend (annual ÷ 12) with deductions.

```protobuf
rpc CalculateMonthlyStipend(CalculateMonthlyStipendRequest) returns (CalculationResponse);
```

**Request:**

```protobuf
message CalculateMonthlyStipendRequest {
  string student_id = 1;
  string stipend_type = 2;
  double annual_amount = 3;
}
```

##### CalculateAnnualStipend

Calculates annual stipend with applicable deductions.

```protobuf
rpc CalculateAnnualStipend(CalculateAnnualStipendRequest) returns (CalculationResponse);
```

**Request:**

```protobuf
message CalculateAnnualStipendRequest {
  string student_id = 1;
  string stipend_type = 2;
  double amount = 3;
}
```

##### CreateStipend

Creates a new stipend record for a student.

```protobuf
rpc CreateStipend(CreateStipendRequest) returns (StipendResponse);
```

**Request:**

```protobuf
message CreateStipendRequest {
  string student_id = 1;
  string stipend_type = 2;
  double amount = 3;
  string payment_method = 4;
  string journal_number = 5;
  string notes = 6;
}
```

**Response:**

```protobuf
message StipendResponse {
  string id = 1;
  string student_id = 2;
  double amount = 3;
  string stipend_type = 4;
  int64 payment_date = 5; // Unix timestamp
  string payment_status = 6;
  string payment_method = 7;
  string journal_number = 8;
  string notes = 9;
  int64 created_at = 10;
  int64 modified_at = 11;
}
```

##### GetStipend

Retrieves a stipend by ID.

```protobuf
rpc GetStipend(GetStipendRequest) returns (StipendResponse);
```

##### GetStudentStipends

Retrieves all stipends for a student with pagination.

```protobuf
rpc GetStudentStipends(GetStudentStipendsRequest) returns (StudentStipendsResponse);
```

**Request:**

```protobuf
message GetStudentStipendsRequest {
  string student_id = 1;
  int32 limit = 2;
  int32 offset = 3;
}
```

**Response:**

```protobuf
message StudentStipendsResponse {
  repeated StipendResponse stipends = 1;
  int32 total = 2;
}
```

##### UpdateStipendPaymentStatus

Updates the payment status of a stipend.

```protobuf
rpc UpdateStipendPaymentStatus(UpdateStipendPaymentStatusRequest) returns (StipendResponse);
```

**Request:**

```protobuf
message UpdateStipendPaymentStatusRequest {
  string stipend_id = 1;
  string payment_status = 2; // "Pending", "Processed", "Failed"
  string payment_date = 3; // ISO8601 format
}
```

---

### 2. DeductionService

The DeductionService provides RPCs for managing deduction rules and applying deductions.

#### RPCs

##### CreateDeductionRule

Creates a new deduction rule.

```protobuf
rpc CreateDeductionRule(CreateDeductionRuleRequest) returns (DeductionRuleResponse);
```

**Request:**

```protobuf
message CreateDeductionRuleRequest {
  string rule_name = 1;
  string deduction_type = 2;
  double default_amount = 3;
  double min_amount = 4;
  double max_amount = 5;
  string frequency = 6; // "Monthly", "Annual", "OneTime"
  bool is_mandatory = 7;
  string applicable_to = 8; // "All", "FullScholarship", "SelfFunded"
  int32 priority = 9;
  string description = 10;
}
```

**Response:**

```protobuf
message DeductionRuleResponse {
  string id = 1;
  string rule_name = 2;
  string deduction_type = 3;
  double default_amount = 4;
  double min_amount = 5;
  double max_amount = 6;
  string frequency = 7;
  bool is_mandatory = 8;
  string applicable_to = 9;
  int32 priority = 10;
  string description = 11;
  bool is_active = 12;
  int64 created_at = 13;
  int64 modified_at = 14;
}
```

##### GetDeductionRule

Retrieves a deduction rule by ID.

```protobuf
rpc GetDeductionRule(GetDeductionRuleRequest) returns (DeductionRuleResponse);
```

##### ListDeductionRules

Lists all deduction rules with optional filters.

```protobuf
rpc ListDeductionRules(ListDeductionRulesRequest) returns (DeductionRulesResponse);
```

**Request:**

```protobuf
message ListDeductionRulesRequest {
  string applicable_to = 1; // Optional filter
  bool is_active = 2;
  int32 limit = 3;
  int32 offset = 4;
}
```

**Response:**

```protobuf
message DeductionRulesResponse {
  repeated DeductionRuleResponse rules = 1;
  int32 total = 2;
}
```

##### ApplyDeductions

Applies deductions to a stipend.

```protobuf
rpc ApplyDeductions(ApplyDeductionsRequest) returns (DeductionsResponse);
```

**Request:**

```protobuf
message ApplyDeductionsRequest {
  string student_id = 1;
  string stipend_id = 2;
  string stipend_type = 3;
  double base_amount = 4;
  repeated string deduction_rule_ids = 5; // Optional: specific rules to apply
}
```

**Response:**

```protobuf
message DeductionsResponse {
  repeated DeductionResponse deductions = 1;
  double total_amount = 2;
  int32 total = 3;
}
```

##### CreateDeduction

Creates a new deduction record.

```protobuf
rpc CreateDeduction(CreateDeductionRequest) returns (DeductionResponse);
```

**Request:**

```protobuf
message CreateDeductionRequest {
  string student_id = 1;
  string deduction_rule_id = 2;
  string stipend_id = 3;
  double amount = 4;
  string deduction_type = 5;
  string description = 6;
}
```

**Response:**

```protobuf
message DeductionResponse {
  string id = 1;
  string student_id = 2;
  string deduction_rule_id = 3;
  string stipend_id = 4;
  double amount = 5;
  string deduction_type = 6;
  string description = 7;
  int64 deduction_date = 8;
  string processing_status = 9;
  string approved_by = 10;
  int64 approval_date = 11;
  string rejection_reason = 12;
  string transaction_id = 13;
  int64 created_at = 14;
  int64 modified_at = 15;
}
```

##### GetDeduction

Retrieves a deduction by ID.

```protobuf
rpc GetDeduction(GetDeductionRequest) returns (DeductionResponse);
```

##### GetStipendDeductions

Retrieves all deductions for a stipend.

```protobuf
rpc GetStipendDeductions(GetStipendDeductionsRequest) returns (DeductionsResponse);
```

##### GetStudentDeductions

Retrieves all deductions for a student.

```protobuf
rpc GetStudentDeductions(GetStudentDeductionsRequest) returns (DeductionsResponse);
```

## Testing

### Running Tests

To run the integration tests for gRPC services:

```bash
# Start the gRPC server first
go run main.go &

# Run the tests
go test -v ./internal/grpc/
```

### Test Files

- `internal/grpc/stipend_server_test.go` - Tests for StipendService
- `internal/grpc/deduction_server_test.go` - Tests for DeductionService

### Test Coverage

Tests cover:

- Stipend calculation with and without deductions
- Monthly and annual stipend calculations
- Creating, retrieving, and updating stipends
- Creating and listing deduction rules
- Creating and retrieving deductions
- Applying deductions to stipends
- Student-specific stipend and deduction queries

## Usage Examples

### Go Client Example

```go
package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "finance_service/pkg/pb"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create client
	client := pb.NewStipendServiceClient(conn)

	// Calculate stipend with deductions
	ctx := context.Background()
	resp, err := client.CalculateStipendWithDeductions(ctx, &pb.CalculateStipendRequest{
		StudentId:     uuid.New().String(),
		StipendType:   "full-scholarship",
		Amount:        100000.00,
		PaymentMethod: "Bank_transfer",
		JournalNumber: "JN-001-2024",
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Base: %.2f, Deductions: %.2f, Net: %.2f",
		resp.BaseStipendAmount, resp.TotalDeductions, resp.NetStipendAmount)
}
```

### Python Client Example

Install the gRPC Python packages:

```bash
pip install grpcio grpcio-tools protobuf
```

Generate Python code:

```bash
python -m grpc_tools.protoc -I proto --python_out=. --grpc_python_out=. proto/stipend.proto proto/deduction.proto
```

Use the client:

```python
import grpc
from finance_service import stipend_pb2, stipend_pb2_grpc
import uuid

def calculate_stipend():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = stipend_pb2_grpc.StipendServiceStub(channel)

        request = stipend_pb2.CalculateStipendRequest(
            student_id=str(uuid.uuid4()),
            stipend_type='full-scholarship',
            amount=100000.00,
            payment_method='Bank_transfer',
            journal_number='JN-001-2024'
        )

        response = stub.CalculateStipendWithDeductions(request)
        print(f"Net Stipend: {response.net_stipend_amount}")

if __name__ == '__main__':
    calculate_stipend()
```

## Proto Files

Proto files are located in the `proto/` directory:

- `proto/stipend.proto` - StipendService definition
- `proto/deduction.proto` - DeductionService definition

## Generated Code

Generated gRPC code is in `pkg/pb/`:

- `pkg/pb/stipend.pb.go` - Stipend message definitions
- `pkg/pb/stipend_grpc.pb.go` - Stipend service stubs
- `pkg/pb/deduction.pb.go` - Deduction message definitions
- `pkg/pb/deduction_grpc.pb.go` - Deduction service stubs

## Error Handling

All gRPC methods return appropriate gRPC status codes:

- `codes.OK` - Success
- `codes.InvalidArgument` - Invalid input
- `codes.NotFound` - Resource not found
- `codes.Internal` - Server error

Errors include descriptive messages for debugging.

## Performance Considerations

- All calculations are performed in-memory
- Database queries use appropriate indexes
- Pagination is supported for list operations
- Deductions are applied in priority order

## Future Enhancements

- Add authentication/authorization
- Implement caching for deduction rules
- Add streaming APIs for batch calculations
- Support for custom deduction formulas
- Metrics and monitoring integration
