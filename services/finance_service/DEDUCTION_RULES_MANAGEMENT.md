# Deduction Rules Management System

## Overview

The Deduction Rules Management System provides comprehensive CRUD (Create, Read, Update, Delete, List) operations for managing deduction rules in the Finance Service. The system is built with proper validation, error handling, and audit trails.

## Features

### 1. **Create Deduction Rule**

- Create new deduction rules with full validation
- Automatic duplicate rule name detection
- Track creator via `created_by` field
- Ensure business logic constraints (min/max amounts, frequency requirements)

### 2. **Read Operations**

- `GetRuleByID`: Fetch a specific rule by its UUID
- `ListActiveRules`: Paginated list of active rules only
- `ListAllRules`: Paginated list including inactive rules
- `ListRulesByType`: Filter rules by deduction type
- `GetApplicableRules`: Get rules applicable to specific student type (full-scholar/self-funded)

### 3. **Update Deduction Rule**

- Partial updates (update only specified fields)
- Comprehensive validation of new values
- Duplicate rule name detection (excluding current rule)
- Audit trail via `modified_by` and `modified_at` fields
- Preserves existing values for unspecified fields

### 4. **Delete Deduction Rule**

- Soft delete (marks rule as inactive instead of hard delete)
- Maintains referential integrity with existing deductions
- Preserves audit trail
- Deleted rules no longer appear in active rule lists

### 5. **Error Logging & Alerts**

- Integration with ErrorLogger for validation and database errors
- Categorized logging:
  - `CategoryDeductionValidation`: Validation errors
  - `CategoryDatabaseError`: Database operation errors
- Automatic logging of all CRUD operations
- Error alerts triggered on thresholds

## Data Model

```go
type DeductionRule struct {
    ID                        uuid.UUID  // Primary key
    RuleName                  string     // Unique rule name (max 100 chars)
    DeductionType             string     // Type: hostel, electricity, mess, etc.
    Description               string     // Rule description
    BaseAmount                float64    // Base deduction amount
    MaxDeductionAmount        float64    // Maximum allowed deduction
    MinDeductionAmount        float64    // Minimum deduction (default 0)
    IsApplicableToFullScholar bool       // Applicable to full scholarships?
    IsApplicableToSelfFunded  bool       // Applicable to self-funded?
    IsActive                  bool       // Is rule currently active?
    AppliesMonthly            bool       // Applied monthly?
    AppliesAnnually           bool       // Applied annually?
    IsOptional                bool       // Optional or mandatory?
    Priority                  int        // Execution priority (higher = first)
    CreatedBy                 *uuid.UUID // User who created the rule
    CreatedAt                 time.Time  // Creation timestamp
    ModifiedBy                *uuid.UUID // User who last modified
    ModifiedAt                time.Time  // Last modification timestamp
}
```

## API Endpoints

### Create Rule

```
POST /api/deduction-rules
Content-Type: application/json

{
  "rule_name": "Hostel Fee",
  "deduction_type": "hostel",
  "description": "Monthly hostel fee",
  "base_amount": 5000.00,
  "max_deduction_amount": 5000.00,
  "min_deduction_amount": 0.00,
  "is_applicable_to_full_scholar": false,
  "is_applicable_to_self_funded": true,
  "applies_monthly": true,
  "applies_annually": false,
  "is_optional": false,
  "priority": 1
}

Response: 201 Created
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "rule_name": "Hostel Fee",
  ...
  "is_active": true,
  "created_at": "2025-11-25T06:00:00Z"
}
```

### Get Rule by ID

```
GET /api/deduction-rules/{ruleID}

Response: 200 OK
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "rule_name": "Hostel Fee",
  ...
}
```

### List Active Rules

```
GET /api/deduction-rules?limit=20&offset=0

Response: 200 OK
{
  "rules": [
    { /* rule 1 */ },
    { /* rule 2 */ }
  ],
  "total": 50,
  "limit": 20,
  "offset": 0
}
```

### List Rules by Type

```
GET /api/deduction-rules/type/{deductionType}?limit=20&offset=0

Response: 200 OK
{
  "rules": [ /* hostel rules only */ ],
  "total": 5,
  "limit": 20,
  "offset": 0
}
```

### Update Rule

```
PUT /api/deduction-rules/{ruleID}
Content-Type: application/json

{
  "rule_name": "Updated Hostel Fee",
  "base_amount": 5500.00
}

Response: 200 OK
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "rule_name": "Updated Hostel Fee",
  "base_amount": 5500.00,
  ...
  "modified_at": "2025-11-25T07:00:00Z"
}
```

### Delete Rule (Soft Delete)

```
DELETE /api/deduction-rules/{ruleID}

Response: 204 No Content
```

## Service Usage Examples

### Create a Rule

```go
service := NewDeductionRuleService()

input := &CreateDeductionRuleInput{
    RuleName:                  "Electricity Bill",
    DeductionType:             "electricity",
    Description:               "Monthly electricity charges",
    BaseAmount:                1500.00,
    MaxDeductionAmount:        2000.00,
    MinDeductionAmount:        500.00,
    IsApplicableToFullScholar: true,
    IsApplicableToSelfFunded:  true,
    AppliesMonthly:            true,
    AppliesAnnually:           false,
    IsOptional:                false,
    Priority:                  2,
    CreatedBy:                 &userID,
}

rule, err := service.CreateRule(input)
if err != nil {
    log.Printf("Error creating rule: %v", err)
}
```

### Update a Rule

```go
input := &UpdateDeductionRuleInput{
    BaseAmount: &newAmount,
    Priority:   &newPriority,
    ModifiedBy: &userID,
}

updatedRule, err := service.UpdateRule(ruleID, input)
if err != nil {
    log.Printf("Error updating rule: %v", err)
}
```

### Delete a Rule

```go
err := service.DeleteRule(ruleID)
if err != nil {
    log.Printf("Error deleting rule: %v", err)
}
```

### List Applicable Rules

```go
// Get rules for a self-funded student
rules, err := service.GetApplicableRules(false)
if err != nil {
    log.Printf("Error fetching applicable rules: %v", err)
}

// Apply rules in priority order
sort.Slice(rules, func(i, j int) bool {
    return rules[i].Priority > rules[j].Priority
})
```

## Validation Rules

### Create/Update Validation

1. **Rule Name**

   - Required, non-empty
   - Maximum 100 characters
   - Must be unique across all rules

2. **Deduction Type**

   - Required, non-empty
   - Examples: "hostel", "electricity", "mess", "other"

3. **Amounts**

   - All amounts must be non-negative
   - MaxDeductionAmount ≥ MinDeductionAmount
   - BaseAmount, MaxDeductionAmount, MinDeductionAmount must be ≥ 0

4. **Frequency**

   - Rule must apply either monthly OR annually (or both)
   - Cannot have both false

5. **Student Type Applicability**
   - At least one of IsApplicableToFullScholar or IsApplicableToSelfFunded must be true

## Error Handling

### Common Errors

| Error                    | HTTP Status               | Example                                               |
| ------------------------ | ------------------------- | ----------------------------------------------------- |
| Rule name already exists | 400 Bad Request           | `"rule name 'Hostel Fee' already exists"`             |
| Rule not found           | 404 Not Found             | `"deduction rule not found"`                          |
| Invalid amount           | 400 Bad Request           | `"base amount cannot be negative"`                    |
| Validation failed        | 400 Bad Request           | `"rule must apply either monthly or annually"`        |
| Database error           | 500 Internal Server Error | `"failed to create deduction rule: connection error"` |

## Error Logging

All operations are logged with appropriate categories:

```
INFO: Deduction rule created: Hostel Fee (rule_id: 550e8400...)
ERROR: Deduction rule creation validation failed (error: base amount cannot be negative)
ERROR: Failed to create deduction rule (rule_name: Test, error: unique constraint violation)
```

## Performance Considerations

- **Indexing**: Rules are indexed by `is_active`, `deduction_type`, and `rule_name`
- **Pagination**: Always use limit/offset for large result sets
- **Priority Sorting**: Rules are sorted by priority (DESC) then name (ASC) for consistent application order
- **Soft Deletes**: Queries automatically filter out inactive rules

## Database Schema

```sql
CREATE TABLE deduction_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rule_name VARCHAR(100) NOT NULL UNIQUE,
    deduction_type VARCHAR(100) NOT NULL,
    description TEXT,
    base_amount DECIMAL(10,2) NOT NULL CHECK (base_amount >= 0),
    max_deduction_amount DECIMAL(10,2) NOT NULL CHECK (max_deduction_amount >= 0),
    min_deduction_amount DECIMAL(10,2) DEFAULT 0 CHECK (min_deduction_amount >= 0),
    is_applicable_to_full_scholar BOOLEAN DEFAULT FALSE,
    is_applicable_to_self_funded BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    applies_monthly BOOLEAN DEFAULT FALSE,
    applies_annually BOOLEAN DEFAULT FALSE,
    is_optional BOOLEAN DEFAULT FALSE,
    priority INTEGER DEFAULT 0,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_by UUID REFERENCES users(id),
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_is_active (is_active),
    INDEX idx_deduction_type (deduction_type),
    INDEX idx_rule_name (rule_name)
);
```

## Testing

Run tests:

```bash
go test ./services -v -run "TestCreateDeductionRule|TestValidationErrors|TestUpdateDeductionRule|TestDeleteDeductionRule"
```

## Future Enhancements

- [ ] Bulk import/export of rules (CSV)
- [ ] Rule templates for common scenarios
- [ ] Audit trail API to view change history
- [ ] Rule versioning system
- [ ] Conditional rules based on student attributes
- [ ] Rule scheduling (active/inactive dates)
- [ ] Integration with payment processing
