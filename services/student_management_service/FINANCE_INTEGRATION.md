# Finance Service Integration

## Overview

The Student Management Service is now fully integrated with the Finance Service to handle stipend calculations, deductions, and financial record management.

## Integration Features

### 1. **Automatic Stipend Calculation**

When creating a stipend allocation, the service:

- Validates student eligibility
- Calls the Finance Service to calculate net stipend amount with deductions
- Applies automatic deductions based on student type (scholarship/self-financed)
- Stores the calculated net amount in the allocation record

### 2. **Finance Service Synchronization**

When creating stipend history records, the service:

- Automatically creates corresponding records in the Finance Service
- Maintains data consistency across both services
- Logs warnings if synchronization fails (graceful degradation)

### 3. **New Integration Endpoints**

#### Calculate Stipend with Deductions

**POST** `/api/stipend/calculate`

Calculates stipend amount with automatic deductions via Finance Service.

**Request Body:**

```json
{
  "student_id": 123,
  "stipend_type": "scholarship",
  "amount": 5000.0
}
```

**Response:**

```json
{
  "base_stipend_amount": 5000.0,
  "total_deductions": 500.0,
  "net_stipend_amount": 4500.0,
  "deductions": [
    {
      "rule_id": "uuid",
      "rule_name": "Health Insurance",
      "deduction_type": "percentage",
      "amount": 250.0,
      "description": "5% health insurance",
      "is_optional": false
    }
  ]
}
```

#### Get Student Finance Stipends

**GET** `/api/students/{studentId}/finance-stipends`

Retrieves all stipend records from the Finance Service for a specific student.

**Response:**

```json
{
  "stipends": [
    {
      "id": "stipend-uuid",
      "student_id": "123",
      "amount": 4500.0,
      "stipend_type": "scholarship",
      "payment_date": 1234567890,
      "payment_status": "pending",
      "payment_method": "bank_transfer",
      "journal_number": "JRN-2025-001",
      "notes": "Monthly stipend",
      "created_at": 1234567890,
      "modified_at": 1234567890
    }
  ],
  "total": 1
}
```

## Configuration

### Environment Variables

Ensure the following environment variable is set:

```bash
FINANCE_GRPC_URL=finance_services:50055
```

In Docker Compose, this is automatically configured:

```yaml
student_management_service:
  environment:
    FINANCE_GRPC_URL: "finance_services:50055"
```

### Service Dependencies

The Student Management Service now depends on:

1. PostgreSQL Database
2. Finance Service (gRPC on port 50055)

## Error Handling

The integration uses **graceful degradation**:

- If Finance Service is unavailable, stipend allocations still work (without deduction calculations)
- Warnings are logged for failed finance service calls
- The service continues to operate independently

## Data Flow

### Creating Stipend Allocation

```
Client Request
    ↓
Student Management Service
    ↓ (validates eligibility)
    ├→ Finance Service (calculate deductions) → returns net amount
    ↓
Saves allocation with net amount
    ↓
Returns response to client
```

### Creating Stipend History

```
Client Request
    ↓
Student Management Service
    ↓
    ├→ Finance Service (create stipend record)
    ├→ Local DB (create history record)
    ↓
Returns response to client
```

## Testing the Integration

### 1. Test Stipend Calculation

```bash
curl -X POST http://localhost:8084/api/stipend/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "stipend_type": "scholarship",
    "amount": 5000.00
  }'
```

### 2. Create Stipend Allocation (with auto-calculation)

```bash
curl -X POST http://localhost:8084/api/stipend/allocations \
  -H "Content-Type: application/json" \
  -d '{
    "allocation_id": "ALLOC-2025-001",
    "student_id": 1,
    "amount": 5000.00,
    "semester": 1,
    "academic_year": "2025"
  }'
```

### 3. Get Finance Stipends for Student

```bash
curl http://localhost:8084/api/students/1/finance-stipends
```

## Benefits

1. **Centralized Financial Logic**: All deduction rules and calculations are managed in the Finance Service
2. **Consistency**: Student stipend data is synchronized across both services
3. **Flexibility**: Finance Service can be updated independently without affecting student management
4. **Auditability**: Complete financial records are maintained in the Finance Service
5. **Scalability**: Services can scale independently based on load

## Architecture

```
┌─────────────────────────────────┐
│  Student Management Service     │
│  (Port 8084, gRPC 50054)       │
│                                 │
│  - Student Records              │
│  - Eligibility Checks           │
│  - Stipend Allocations          │
│  - Stipend History              │
└────────────┬────────────────────┘
             │
             │ gRPC (port 50055)
             │
             ↓
┌─────────────────────────────────┐
│      Finance Service            │
│  (Port 8085, gRPC 50055)       │
│                                 │
│  - Stipend Calculations         │
│  - Deduction Rules              │
│  - Financial Records            │
│  - Transaction Management       │
└─────────────────────────────────┘
```

## Next Steps

1. Test the integration with real student data
2. Configure deduction rules in the Finance Service
3. Set up monitoring for gRPC communication
4. Implement retry logic for critical operations
5. Add circuit breaker pattern for resilience
