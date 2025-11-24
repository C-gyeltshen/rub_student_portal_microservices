# COMPLETE SERVICE INTEGRATION & DATABASE FLOW DIAGRAM

SYSTEM ARCHITECTURE DIAGRAM
═══════════════════════════

┌─────────────────────────────────────────────────────────────────────┐
│ CLIENT / FRONTEND │
└──────────────────────────────┬──────────────────────────────────────┘
│
▼
┌──────────────────────┐
│ API GATEWAY │
│ (Port 8080) │
│ Routes Requests │
└──────┬───────────────┘
│
┌──────────────────┼──────────────────┐
│ │ │
▼ ▼ ▼
┌────────┐ ┌──────────┐ ┌──────────────┐
│ USER │ │ BANKING │ │ FINANCE │
│SERVICE │ │ SERVICE │ │ SERVICE │
│(8082) │ │ (8083) │ │ (8084) │
└────┬───┘ └────┬─────┘ └──────┬───────┘
│ │ │
└──────────────────┴─────────────────┘
│
▼
┌──────────────────┐
│ PostgreSQL │
│ Database │
│ (Port 5434) │
│ Single Instance │
└──────────────────┘

DATABASE SCHEMA & RELATIONSHIPS
═══════════════════════════════

┌──────────────────────────────────────────────────────────────────────┐
│ SHARED DATABASE │
│ (rub_student_portal - PostgreSQL) │
└──────────────────────────────────────────────────────────────────────┘

┌─────────────────────┐ ┌──────────────────┐ ┌──────────────┐
│ USERS │ │ ROLES │ │ COLLEGES │
├─────────────────────┤ ├──────────────────┤ ├──────────────┤
│ id (PK) │ │ id (PK) │ │ id (PK) │
│ email │ │ name (UNIQUE) │ │ name (UNIQUE)│
│ hash_password │ │ created_at │ │ created_at │
│ role_id (FK)────────┼─→│ modified_at │ │ modified_at │
│ created_at │ └──────────────────┘ └──────────────┘
│ modified_at │
└─────────────────────┘

┌──────────────────────┐ ┌──────────────────────┐
│ PROGRAMS │ │ STUDENTS │
├──────────────────────┤ ├──────────────────────┤
│ id (PK) │ │ id (PK) │
│ name (UNIQUE) │ │ user_id (FK) [UNIQ] │
│ college_id (FK)──────┼─→│ name │
│ modified_by (FK) │ │ rub_id_card_number │
│ created_at │ │ email │
│ modified_at │ │ phone_number │
└──────────────────────┘ │ date_of_birth │
│ college_id (FK) │
│ program_id (FK) │
│ modified_by (FK) │
│ created_at │
│ modified_at │
└──────────────────────┘

┌──────────────────┐ ┌─────────────────────────────┐
│ BANKS │ │ STUDENT_BANK_DETAILS │
├──────────────────┤ ├─────────────────────────────┤
│ id (PK) │ │ id (PK) │
│ name (UNIQUE) │ │ student_id (FK) [UNIQUE] │
│ created_at │ │ bank_id (FK)────────────────→
│ modified_at │ │ account_number (UNIQUE) │
└──────────────────┘ │ account_holder_name │
│ created_at │
│ modified_at │
└─────────────────────────────┘

FINANCE TABLES:
═══════════════

┌─────────────────────────────────────┐
│ STIPENDS │
├─────────────────────────────────────┤
│ id (PK) │
│ student_id (FK)─────────────────→ │ (to STUDENTS)
│ amount (DECIMAL 10,2) [>= 0] │
│ stipend_type VARCHAR(50) │
│ payment_date │
│ payment_status VARCHAR(50) │
│ payment_method VARCHAR(50) │
│ journal_number TEXT (UNIQUE) │
│ transaction_id (FK) ──────────→ │ (to TRANSACTION)
│ notes TEXT │
│ created_at │
│ modified_at │
│ INDEX: student_id, status, type │
└─────────────────────────────────────┘

┌──────────────────────────────────────┐
│ DEDUCTION_RULES │
├──────────────────────────────────────┤
│ id (PK) │
│ rule_name VARCHAR(100) (UNIQUE) │
│ deduction_type VARCHAR(100) │
│ description TEXT │
│ base_amount DECIMAL(10,2) [>= 0] │
│ max_deduction DECIMAL(10,2) [>= 0] │
│ min_deduction DECIMAL(10,2) [>= 0] │
│ is_applicable_to_full_scholar BOOL │
│ is_applicable_to_self_funded BOOL │
│ is_active BOOL (INDEX) │
│ applies_monthly BOOL │
│ applies_annually BOOL │
│ is_optional BOOL │
│ priority INTEGER │
│ created_by (FK)──────────────────→ │ (to USERS)
│ modified_by (FK)─────────────────→ │ (to USERS)
│ created_at │
│ modified_at │
│ INDEX: rule_name, is_active, type │
└──────────────────────────────────────┘

┌──────────────────────────────────────────┐
│ DEDUCTIONS │
├──────────────────────────────────────────┤
│ id (PK) │
│ student_id (FK)──────────────────────→ │ (to STUDENTS)
│ deduction_rule_id (FK) [RESTRICT]────→ │ (to DEDUCTION_RULES)
│ stipend_id (FK)────────────────────────→ │ (to STIPENDS)
│ amount DECIMAL(10,2) [>= 0] │
│ deduction_type VARCHAR(100) │
│ description TEXT │
│ deduction_date │
│ processing_status VARCHAR(50) │
│ approved_by (FK)──────────────────────→ │ (to USERS)
│ approval_date │
│ rejection_reason TEXT │
│ transaction_id (FK) ─────────────────→ │ (to TRANSACTION)
│ created_at │
│ modified_at │
│ INDEX: student_id, stipend_id, status │
└──────────────────────────────────────────┘

┌────────────────────────────────────┐
│ FINANCE_OFFICER │
├────────────────────────────────────┤
│ id (PK) │
│ officer_user_id (FK) [UNIQUE]──→ │ (to USERS)
│ name TEXT │
│ college_id (FK)────────────────→ │ (to COLLEGES)
│ created_at │
│ modified_at │
└────────────────────────────────────┘

┌─────────────────────────────────────────┐
│ BUDGET │
├─────────────────────────────────────────┤
│ id (PK) │
│ finance_officer_id (FK)─────────────→ │ (to FINANCE_OFFICER)
│ name TEXT │
│ purpose TEXT │
│ allocated_amount DECIMAL(12,2) │
│ start_date DATE │
│ end_date DATE │
│ status VARCHAR(50) │
│ created_at │
│ modified_at │
└─────────────────────────────────────────┘

┌────────────────────────────────────────┐
│ EXPENSES │
├────────────────────────────────────────┤
│ id (PK) │
│ budget_id (FK)──────────────────────→ │ (to BUDGET)
│ requester_id (FK)───────────────────→ │ (to USERS)
│ description TEXT │
│ amount DECIMAL(10,2) │
│ status VARCHAR(50) │
│ approval_date │
│ created_at │
│ modified_at │
└────────────────────────────────────────┘

┌────────────────────────────────────────────┐
│ TRANSACTION │
├────────────────────────────────────────────┤
│ id (PK) │
│ expenses_id (FK) [nullable]──────────────→ │ (to EXPENSES)
│ transaction_type VARCHAR(50) │
│ amount DECIMAL(12,2) │
│ transaction_date │
│ journal_number TEXT (UNIQUE) │
│ bank_id (FK)───────────────────────────→ │ (to BANKS)
│ notes TEXT │
│ created_at │
│ modified_at │
└────────────────────────────────────────────┘

STIPEND CREATION DATA FLOW
═══════════════════════════

┌─────────────────────────────────────────────────────────────────┐
│ REQUEST: POST /api/finance/stipends │
│ BODY: { │
│ "student_id": "uuid-123", │
│ "amount": 50000, │
│ "created_by": "user-456" │
│ } │
└────────────────────────────────┬────────────────────────────────┘
│
┌────────────▼────────────┐
│ Finance Service │
│ Validation Pipeline │
└────────────┬────────────┘
│
┌───────────────────────┼───────────────────────┐
│ │ │
▼ ▼ ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│ UserClient │ │StudentClient│ │BankingClient│
│ VALIDATE │ │ VALIDATE │ │ VALIDATE │
│ USER ROLE │ │ STUDENT │ │ BANK ACCT │
└─────────────┘ └─────────────┘ └─────────────┘
│ │ │
│ Call to User Service │ Call to User Service │ Call to Banking Service
│ GET /api/users/456 │ GET /api/students/123│ GET /api/student-bank/123
└───────────────────────┴───────────────────────┘
│
┌────────────▼────────────┐
│ All Validations Pass │
│ OR Return Error │
└────────────┬────────────┘
│
┌───────────────────────┘
│
▼
┌──────────────────────────────────────┐
│ DATABASE OPERATIONS │
├──────────────────────────────────────┤
│ 1. INSERT into stipends │
│ - student_id, amount, status │
│ - payment_status = 'Pending' │
│ - journal_number = unique gen │
│ │
│ 2. SELECT \* from deduction_rules │
│ WHERE is_active = true │
│ AND applicable to student_type │
│ ORDER BY priority DESC │
│ │
│ 3. INSERT into deductions (loop) │
│ FOR EACH deduction_rule: │
│ - rule_id, stipend_id, amount │
│ - status = 'Pending' │
│ │
│ 4. INSERT into transaction (audit) │
│ - journal_number, user_id │
│ - timestamp, type │
└────────────┬───────────────────────┘
│
▼
┌──────────────────────────────────┐
│ RESPONSE: 200 OK │
│ { │
│ "stipend_id": "uuid-789", │
│ "amount": 50000, │
│ "net_amount": 43000, │
│ "journal_number": "STIP-001", │
│ "deductions": [ │
│ { │
│ "rule_name": "hostel", │
│ "amount": 5000 │
│ }, │
│ { │
│ "rule_name": "mess", │
│ "amount": 2000 │
│ } │
│ ], │
│ "created_at": "2025-11-25..." │
│ } │
└──────────────────────────────────┘

SERVICE COMMUNICATION MATRIX
════════════════════════════

┌──────────────────┬──────────────┬─────────────────────────────┐
│ Finance Service │ Target │ Purpose │
├──────────────────┼──────────────┼─────────────────────────────┤
│ StudentClient │ User Svc: │ Validate student exists, │
│ │ 8082 │ fetch info, check account │
├──────────────────┼──────────────┼─────────────────────────────┤
│ BankingClient │ Banking Svc: │ Validate bank details, │
│ │ 8083 │ fetch account info │
├──────────────────┼──────────────┼─────────────────────────────┤
│ UserClient │ User Svc: │ Validate permissions, │
│ │ 8082 │ fetch user role, audit │
├──────────────────┼──────────────┼─────────────────────────────┤
│ Database Ops │ PostgreSQL: │ CRUD on stipends, │
│ │ 5434 │ deductions, transactions │
└──────────────────┴──────────────┴─────────────────────────────┘

ERROR HANDLING FLOW
═══════════════════

┌─────────────────────────────────────┐
│ Validation Error Detected │
└────────────────────┬────────────────┘
│
┌───────────┼───────────┐
│ │ │
▼ ▼ ▼
┌─────────┐ ┌────────┐ ┌──────────┐
│ Service │ │Amount │ │Deduction │
│ Error │ │Error │ │Error │
└────┬────┘ └────┬───┘ └────┬─────┘
│ │ │
▼ ▼ ▼
ROLLBACK: No database changes
│
▼
RETURN 400/401:
{
"error": "Clear error message",
"details": "Specific validation failure"
}

DEPLOYMENT TOPOLOGY
═══════════════════

┌─────────────────────────────────────────────────────┐
│ Docker Compose Network: rub-network │
├─────────────────────────────────────────────────────┤
│ │
│ ┌──────────────────────────────────────────────┐ │
│ │ API Gateway Container │ │
│ │ - Image: golang:1.23-alpine (build) │ │
│ │ - Ports: 8080:8080 │ │
│ │ - Routes: finance, user, banking requests │ │
│ └──────────────────────────────────────────────┘ │
│ │ │
│ ┌─────────────────┼─────────────────┐ │
│ │ │ │ │
│ ┌─▼──────────┐ ┌───▼──────┐ ┌────────▼─┐ │
│ │ Finance │ │ User │ │ Banking │ │
│ │ Service │ │ Service │ │ Service │ │
│ │ :8084 │ │ :8082 │ │ :8083 │ │
│ │ (golang) │ │ (golang) │ │ (golang) │ │
│ └─┬──────────┘ └───┬──────┘ └────────┬─┘ │
│ │ │ │ │
│ └────────────────┼─────────────────┘ │
│ │ │
│ ┌─────▼──────┐ │
│ │ PostgreSQL │ │
│ │ Container │ │
│ │ Port: 5434 │ │
│ │ Database │ │
│ │ Volume │ │
│ └────────────┘ │
│ │
└─────────────────────────────────────────────────────┘

This completes the comprehensive service integration architecture.
All components are ready for 1.2 Stipend Calculation Logic implementation.
