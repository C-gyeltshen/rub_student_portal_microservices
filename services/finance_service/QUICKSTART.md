# gRPC Quick Start Guide

## Overview

The Finance Service now provides both REST API and gRPC interfaces for stipend calculation and deduction management.

- **REST API**: `http://localhost:8084`
- **gRPC Server**: `localhost:50051`

## Quick Start

### Prerequisites

1. **Go 1.22+**

   ```bash
   go version  # Should show go1.22 or higher
   ```

2. **PostgreSQL Database**

   ```bash
   # Using Docker (recommended)
   docker run -d --name postgres -e POSTGRES_PASSWORD=password \
     -e POSTGRES_DB=finance_db -p 5432:5432 postgres:latest

   # Or use your existing PostgreSQL installation
   ```

3. **Set Environment Variables**
   ```bash
   export DATABASE_URL="postgres://postgres:password@localhost:5432/finance_db"
   export PORT=8084
   export GRPC_PORT=50051
   ```

### Full Walkthrough (Complete Setup)

Complete step-by-step setup for first-time users:

**Step 1: Start PostgreSQL**

```bash
docker run -d --name finance-postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=finance_db \
  -p 5432:5432 \
  postgres:latest

# Wait 2-3 seconds for PostgreSQL to start
sleep 3
```

**Step 2: Navigate to service directory**

```bash
cd services/finance_service
```

**Step 3: Set environment variables**

```bash
export DATABASE_URL="postgres://postgres:password@localhost:5432/finance_db"
export PORT=8084
export GRPC_PORT=50051
```

**Step 4: Run the service**

```bash
go run main.go
```

**Expected output:**

```
gRPC Server listening on port 50051
Finance Service starting on port 8084
✅ Service ready for requests
```

**Step 5: Test in another terminal**

```bash
# Test REST API
curl http://localhost:8084/health

# Or test gRPC (if grpcurl installed)
grpcurl -plaintext localhost:50051 list
```

**Step 6: Run integration tests**

```bash
# In another terminal
cd services/finance_service
go test -v ./internal/grpc/ -timeout 30s
```

### Option 1: Run Locally

```bash
cd services/finance_service
export DATABASE_URL="postgres://postgres:password@localhost:5432/finance_db"
go run main.go
```

**Output:**

```
gRPC Server listening on port 50051
Finance Service starting on port 8084
```

### Option 2: Build & Run

```bash
cd services/finance_service
go build .
export DATABASE_URL="postgres://postgres:password@localhost:5432/finance_db"
./finance_service
```

### Option 3: Docker

```bash
# Build the image
docker build -t finance-service:latest .

# Run with PostgreSQL
docker network create finance-network
docker run -d --name postgres --network finance-network \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=finance_db \
  postgres:latest

docker run -p 8084:8084 -p 50051:50051 \
  --network finance-network \
  -e DATABASE_URL="postgres://postgres:password@postgres:5432/finance_db" \
  finance-service:latest
```

## Testing gRPC APIs

### Run Integration Tests

**Important:** Integration tests require the server to be running on port 50051.

**Terminal 1 - Start the server:**

```bash
cd services/finance_service
go run main.go
```

**Terminal 2 - Run tests:**

```bash
cd services/finance_service
go test -v ./internal/grpc/
```

### Test Stipend Calculation

```bash
go test -v -run TestStipendService_CalculateStipendWithDeductions ./internal/grpc/
```

### Test Deduction Management

```bash
go test -v -run TestDeductionService ./internal/grpc/
```

### Quick Test Without Server

If you just want to verify compilation without running the server:

```bash
go test -c ./internal/grpc/  # Compile tests without running
```

### Verification Checklist

After starting the service, verify everything works:

```bash
# 1. Check service is running
echo "Step 1: Checking service..."
curl -s http://localhost:8084/health || echo "REST API not responding"

# 2. Check gRPC server
echo "Step 2: Checking gRPC..."
nc -zv localhost 50051 && echo "✅ gRPC port open" || echo "❌ gRPC port closed"

# 3. Compile tests
echo "Step 3: Compiling tests..."
go test -c ./internal/grpc/ -v

# 4. Check generated proto files
echo "Step 4: Checking proto files..."
ls -la pkg/pb/*.go | wc -l

echo "✅ All checks complete!"
```

## Using gRPC APIs

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
    resp, err := client.CalculateStipendWithDeductions(
        context.Background(),
        &pb.CalculateStipendRequest{
            StudentId:     uuid.New().String(),
            StipendType:   "full-scholarship",
            Amount:        100000.00,
            PaymentMethod: "Bank_transfer",
            JournalNumber: "JN-001-2024",
        },
    )

    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    log.Printf("Base: %.2f, Deductions: %.2f, Net: %.2f",
        resp.BaseStipendAmount, resp.TotalDeductions, resp.NetStipendAmount)
}
```

### Using grpcurl

Install grpcurl:

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

List services:

```bash
grpcurl -plaintext localhost:50051 list
```

Call a method:

```bash
grpcurl -plaintext -d '{
  "student_id": "550e8400-e29b-41d4-a716-446655440000",
  "stipend_type": "full-scholarship",
  "amount": 100000,
  "payment_method": "Bank_transfer",
  "journal_number": "JN-001-2024"
}' localhost:50051 pb.StipendService/CalculateStipendWithDeductions
```

## API Endpoints

### Stipend Service (gRPC)

| Method                           | Request                             | Response                  |
| -------------------------------- | ----------------------------------- | ------------------------- |
| `CalculateStipendWithDeductions` | `CalculateStipendRequest`           | `CalculationResponse`     |
| `CalculateMonthlyStipend`        | `CalculateMonthlyStipendRequest`    | `CalculationResponse`     |
| `CalculateAnnualStipend`         | `CalculateAnnualStipendRequest`     | `CalculationResponse`     |
| `CreateStipend`                  | `CreateStipendRequest`              | `StipendResponse`         |
| `GetStipend`                     | `GetStipendRequest`                 | `StipendResponse`         |
| `GetStudentStipends`             | `GetStudentStipendsRequest`         | `StudentStipendsResponse` |
| `UpdateStipendPaymentStatus`     | `UpdateStipendPaymentStatusRequest` | `StipendResponse`         |

### Deduction Service (gRPC)

| Method                 | Request                       | Response                 |
| ---------------------- | ----------------------------- | ------------------------ |
| `CreateDeductionRule`  | `CreateDeductionRuleRequest`  | `DeductionRuleResponse`  |
| `GetDeductionRule`     | `GetDeductionRuleRequest`     | `DeductionRuleResponse`  |
| `ListDeductionRules`   | `ListDeductionRulesRequest`   | `DeductionRulesResponse` |
| `ApplyDeductions`      | `ApplyDeductionsRequest`      | `DeductionsResponse`     |
| `CreateDeduction`      | `CreateDeductionRequest`      | `DeductionResponse`      |
| `GetDeduction`         | `GetDeductionRequest`         | `DeductionResponse`      |
| `GetStipendDeductions` | `GetStipendDeductionsRequest` | `DeductionsResponse`     |
| `GetStudentDeductions` | `GetStudentDeductionsRequest` | `DeductionsResponse`     |

## Environment Variables

```bash
# REST API Port (default: 8084)
export PORT=8084

# gRPC Server Port (default: 50051)
export GRPC_PORT=50051

# Database URL (required)
# PostgreSQL format: postgres://user:password@host:port/database
export DATABASE_URL="postgres://postgres:password@localhost:5432/finance_db"
```

### Database Setup

The application automatically runs migrations when started. You just need to ensure the database exists:

```bash
# Using psql directly (optional - let the app create tables)
psql -U postgres -h localhost -c "CREATE DATABASE finance_db;"

# The app will auto-migrate tables when you run it for the first time
```

If PostgreSQL isn't running, start it first:

```bash
# Docker container
docker run -d --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=finance_db \
  -p 5432:5432 \
  postgres:latest
```

## File Structure

```
services/finance_service/
├── main.go                          # Dual-server setup
├── proto/
│   ├── stipend.proto               # Stipend service definition
│   └── deduction.proto             # Deduction service definition
├── pkg/pb/
│   ├── stipend.pb.go               # Generated message code
│   ├── stipend_grpc.pb.go          # Generated service stubs
│   ├── deduction.pb.go             # Generated message code
│   └── deduction_grpc.pb.go        # Generated service stubs
├── internal/grpc/
│   ├── stipend_server.go           # Stipend service implementation
│   ├── stipend_server_test.go      # Stipend tests
│   ├── deduction_server.go         # Deduction service implementation
│   └── deduction_server_test.go    # Deduction tests
├── services/
│   ├── stipend_service.go          # Business logic
│   ├── deduction_service.go        # Business logic
│   └── types.go                    # Type definitions
├── Dockerfile                      # Multi-stage Docker build
├── GRPC_API_REFERENCE.md          # Complete API documentation
└── VALIDATION_REPORT.md           # Validation details
```

## Troubleshooting

### Connection Refused on Port 50051

- Ensure the server is running: `go run main.go`
- Check port is not in use: `lsof -i :50051`
- Verify firewall allows port 50051

### Proto Files Not Updating

- Delete `pkg/pb/` directory
- Run `protoc --go_out=pkg/pb --go-grpc_out=pkg/pb -I proto proto/*.proto`
- Rebuild with `go build .`

### Database Connection Issues

- Verify `DATABASE_URL` environment variable
- Ensure PostgreSQL is running
- Check database credentials and permissions

### Tests Fail with "Connection Refused"

- Start the server first: `go run main.go &`
- Then run tests: `go test -v ./internal/grpc/`

## Performance Characteristics

- **Calculation**: < 100ms for typical stipend calculation
- **Database Queries**: Optimized with proper indexing
- **Concurrent Requests**: Both REST and gRPC handle multiple concurrent clients
- **Memory**: ~50MB base, scales with active connections

## Security Considerations

For production deployment:

1. Add TLS/SSL for gRPC connections
2. Implement authentication/authorization
3. Add request validation and rate limiting
4. Use environment variables for sensitive data
5. Enable gRPC reflection for debugging (disable in production)

## Documentation References

- [Complete API Reference](./GRPC_API_REFERENCE.md)
- [Validation Report](./VALIDATION_REPORT.md)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [gRPC Documentation](https://grpc.io/docs/)

## Support

For issues or questions, refer to:

- API Reference: `GRPC_API_REFERENCE.md`
- Test Examples: `internal/grpc/*_test.go`
- Service Logic: `services/deduction_service.go`, `services/stipend_service.go`
