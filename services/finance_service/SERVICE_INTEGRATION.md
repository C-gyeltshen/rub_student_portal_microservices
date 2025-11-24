# SERVICE INTEGRATION ARCHITECTURE

1. # SERVICE COMMUNICATION FLOW

┌─────────────────────────────────────────────────────────────────┐
│ API GATEWAY │
│ (Port 8080) │
└──────┬─────────────┬──────────────────┬──────────────────────────┘
│ │ │
▼ ▼ ▼
┌────────┐ ┌──────────┐ ┌──────────────┐
│ USER │ │ BANKING │ │ FINANCE │
│SERVICE │ │ SERVICE │ │ SERVICE │
│(8082) │ │ (8083) │ │ (8084) │
└─┬──────┘ └──┬───────┘ └──┬───────────┘
│ │ │
└─────────────┴───────────────┘
│
▼
┌─────────────┐
│ PostgreSQL │
│ (Port 5434)│
└─────────────┘

2. # DATA FLOW FOR STIPEND CREATION

Client Request:
↓
POST /api/finance/stipends
↓
Finance Service
├─ 1. Validate Request (amount, student_id)
│
├─ 2. Call Student Service
│ ├─ GET /api/students/{student_id}
│ ├─ Fetch: name, college, program, stipend_type
│ └─ Validate: student exists and is active
│
├─ 3. Call Banking Service
│ ├─ GET /api/student-bank-details/{student_id}
│ ├─ Fetch: bank_id, account_number, account_holder
│ └─ Validate: banking details exist
│
├─ 4. Fetch User Info (from same DB)
│ ├─ Check role: admin, finance_officer, student
│ └─ Validate: user has permission to create stipend
│
├─ 5. Calculate Applicable Deductions
│ ├─ Query deduction_rules table
│ ├─ Filter by: is_active=true, student_type match
│ └─ Sort by: priority DESC
│
├─ 6. Create Stipend Record
│ ├─ INSERT into stipends table
│ ├─ Set: student_id, amount, payment_status='Pending'
│ └─ Generate: journal_number (unique reference)
│
├─ 7. Create Deduction Records
│ ├─ FOR EACH applicable deduction_rule:
│ │ ├─ Calculate deduction amount
│ │ ├─ INSERT into deductions table
│ │ └─ Set: status='Pending', awaiting_approval
│ │
│ └─ Update: stipend.net_amount = amount - total_deductions
│
├─ 8. Audit Trail
│ ├─ Log: transaction in transaction table
│ ├─ Record: created_by = user_id, created_at = timestamp
│ └─ Link: to expense/budget if applicable
│
└─ 9. Send Response
└─ Return: {stipend_id, amount, deductions[], net_amount}

3. # SHARED DATABASE SCHEMA

SHARED TABLES (All Services Use):
├─ users (from User Service)
├─ roles
├─ students (from Student Service)
├─ colleges
├─ programs
├─ banks (from Banking Service)
├─ student_bank_details (from Banking Service)
└─ Additional for Finance:
├─ stipends
├─ deduction_rules
├─ deductions
├─ finance_officer
├─ budget
├─ expenses
└─ transaction

4. # SERVICE-TO-SERVICE DEPENDENCIES

FINANCE SERVICE DEPENDS ON:
├─ User Service
│ ├─ Purpose: Authentication, Authorization
│ ├─ Data: user_id, role, permissions
│ └─ API: GET /api/users/{id}, POST /api/users/verify-token
│
├─ Student Service  
│ ├─ Purpose: Student information, eligibility check
│ ├─ Data: student_id, name, college_id, program_id, stipend_type
│ └─ API: GET /api/students/{id}, GET /api/students/by-card/{card_id}
│
└─ Banking Service
├─ Purpose: Bank account validation, payment processing
├─ Data: bank_id, account_number, account_holder, routing_number
└─ API: GET /api/banks/{id}, GET /api/student-bank/{student_id}

5. # INTEGRATION POINTS (IN DATABASE)

FOREIGN KEY RELATIONSHIPS:

stipends table:
├─ student_id → students(id)
│ └─ When: Creating stipend for a student
│ Must: Verify student exists in DB
│
├─ transaction_id → transaction(id)
│ └─ When: Stipend is processed/paid
│ Links: to payment transaction record

deductions table:
├─ student_id → students(id)
│ └─ For: Audit trail and student-level tracking
│
├─ deduction_rule_id → deduction_rules(id)
│ └─ Links: Which rule was applied
│
├─ stipend_id → stipends(id)
│ └─ Links: Which stipend this deduction applies to
│
├─ approved_by → users(id)
│ └─ For: Audit - who approved this deduction

deduction_rules table:
├─ created_by → users(id)
├─ modified_by → users(id)
│ └─ For: Audit trail - who created/modified rules

finance_officer table:
├─ officer_user_id → users(id) [UNIQUE]
│ └─ Links: Finance officer to their user account
│
└─ college_id → colleges(id)
└─ Scope: Finance officer manages college budget

budget table:
├─ finance_officer_id → finance_officer(id)
│ └─ Who: Finance officer created this budget

expenses table:
├─ budget_id → budget(id)
│ └─ Which: Budget this expense charges to
│
└─ requester_id → users(id)
└─ Who: Requested/submitted the expense

transaction table:
├─ expenses_id → expenses(id) [nullable]
│ └─ If: This transaction is for an expense
│
├─ bank_id → banks(id)
│ └─ Which: Bank processed this transaction
│
└─ journal_number → UNIQUE
└─ Reference: Audit trail unique identifier

6. # DATA FLOW EXAMPLES

EXAMPLE 1: Create Stipend for Student
─────────────────────────────────────

1. User (role: finance_officer) submits:
   POST /api/finance/stipends
   {
   "student_id": "uuid-123",
   "amount": 50000.00,
   "stipend_type": "full-scholarship"
   }

2. Finance Service executes:
   a) SELECT students WHERE id = 'uuid-123'
   → Fetch student info

   b) SELECT deduction_rules
   WHERE is_active=true
   AND is_applicable_to_full_scholar=true
   ORDER BY priority DESC
   → Get applicable deductions

   c) INSERT into stipends
   → Create stipend record

   d) FOR each deduction_rule:
   INSERT into deductions
   → Create deduction records

   e) INSERT into transaction
   → Create audit record

3. Response:
   {
   "stipend_id": "uuid-456",
   "amount": 50000.00,
   "deductions": [
   {"rule_name": "hostel_fee", "amount": 5000},
   {"rule_name": "mess_fee", "amount": 2000}
   ],
   "net_amount": 43000.00,
   "journal_number": "STIP-2025-001"
   }

EXAMPLE 2: Process Payment
──────────────────────────

1. Stipend status changes to 'Processing'
   UPDATE stipends SET payment_status='Processing'

2. System fetches student bank details
   SELECT \* FROM student_bank_details
   WHERE student_id = 'uuid-123'

3. System initiates payment

   - Call Banking Service API (internal)
   - Send: amount, account_number, bank_id
   - Get: transaction_id, status

4. Update records
   UPDATE stipends SET
   payment_status='Processed',
   transaction_id='txn-789',
   payment_date=NOW()

   UPDATE deductions SET
   processing_status='Processed',
   approval_date=NOW()

5. Audit entry
   INSERT into transaction
   {
   "stipend_id": "uuid-456",
   "type": "DEBIT",
   "amount": 43000.00,
   "journal_number": "STIP-2025-001",
   "bank_id": "bank-123"
   }

6. # ERROR HANDLING & VALIDATION

Student Validation:
├─ Student exists in DB
├─ Student is active (not suspended)
├─ Student has valid enrollment
└─ Student stipend_type is supported

Banking Validation:
├─ Bank account exists for student
├─ Account is active and verified
├─ Account holder name matches student
└─ No duplicate stipends for same period

Financial Validation:
├─ Amount is non-negative
├─ Amount doesn't exceed max allowed
├─ Budget allocation is sufficient
└─ No duplicate journal_number

Permission Validation:
├─ User has finance_officer or admin role
├─ User can manage student's college
└─ User hasn't reached approval limit

Deduction Validation:
├─ Deduction rule is active
├─ Deduction applies to student type
├─ Deduction amount is within limits
├─ Deduction doesn't exceed stipend amount
└─ Total deductions ≤ max_deduction_percent
