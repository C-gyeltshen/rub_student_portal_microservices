-- -----------------------------------------------------------------------
-- Finance Service - Stipend Calculation & Deduction System
-- -----------------------------------------------------------------------

-- Stipends Table
-- Stores stipend payment records for students (full-scholarship or self-funded)
CREATE TABLE IF NOT EXISTS stipends(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount >= 0),
    stipend_type VARCHAR(50) NOT NULL, -- full-scholarship, self-funded
    payment_date TIMESTAMPTZ,
    payment_status VARCHAR(50) DEFAULT 'Pending', -- Pending, Processed, Failed
    payment_method VARCHAR(50), -- Bank_transfer, E-payment
    journal_number TEXT NOT NULL UNIQUE,
    transaction_id UUID REFERENCES transactions(id) ON DELETE SET NULL,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Deduction Rules Table
-- Defines configurable rules for deductions (hostel, electricity, mess fees, etc)
CREATE TABLE IF NOT EXISTS deduction_rules(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rule_name VARCHAR(100) NOT NULL UNIQUE,
    deduction_type VARCHAR(100) NOT NULL, -- hostel, electricity, mess_fees, etc
    description TEXT,
    base_amount DECIMAL(10, 2) NOT NULL CHECK (base_amount >= 0),
    max_deduction_amount DECIMAL(10, 2) NOT NULL CHECK (max_deduction_amount >= 0),
    min_deduction_amount DECIMAL(10, 2) DEFAULT 0 CHECK (min_deduction_amount >= 0),
    is_applicable_to_full_scholar BOOLEAN DEFAULT FALSE,
    is_applicable_to_self_funded BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    applies_monthly BOOLEAN DEFAULT FALSE,
    applies_annually BOOLEAN DEFAULT FALSE,
    is_optional BOOLEAN DEFAULT FALSE, -- true if optional, false if mandatory
    priority INTEGER DEFAULT 0, -- Higher priority deductions applied first
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_by UUID REFERENCES users(id) ON DELETE SET NULL,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Deductions Table
-- Records actual deductions applied to students' stipends
CREATE TABLE IF NOT EXISTS deductions(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    deduction_rule_id UUID NOT NULL REFERENCES deduction_rules(id) ON DELETE RESTRICT,
    stipend_id UUID NOT NULL REFERENCES stipends(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount >= 0),
    deduction_type VARCHAR(100) NOT NULL, -- hostel, electricity, mess_fees, etc
    description TEXT,
    deduction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    processing_status VARCHAR(50) DEFAULT 'Pending', -- Pending, Approved, Processed, Rejected
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approval_date TIMESTAMPTZ,
    rejection_reason TEXT,
    transaction_id UUID REFERENCES transactions(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Transactions Table
-- Records money transfer transactions for stipend distribution
CREATE TABLE IF NOT EXISTS transactions(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stipend_id UUID NOT NULL REFERENCES stipends(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    source_account VARCHAR(255), -- College/Institution account
    destination_account VARCHAR(255) NOT NULL, -- Student's account
    destination_bank VARCHAR(255),
    transaction_type VARCHAR(50) DEFAULT 'STIPEND', -- STIPEND, REFUND, etc
    status VARCHAR(50) DEFAULT 'PENDING', -- PENDING, PROCESSING, SUCCESS, FAILED, CANCELLED
    payment_method VARCHAR(50), -- BANK_TRANSFER, E_PAYMENT, etc
    reference_number VARCHAR(255) UNIQUE, -- Unique reference from payment gateway
    error_message TEXT, -- Error details if transaction failed
    remarks TEXT,
    initiated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_stipends_student_id ON stipends(student_id);
CREATE INDEX IF NOT EXISTS idx_stipends_payment_status ON stipends(payment_status);
CREATE INDEX IF NOT EXISTS idx_stipends_stipend_type ON stipends(stipend_type);
CREATE INDEX IF NOT EXISTS idx_stipends_created_at ON stipends(created_at);

CREATE INDEX IF NOT EXISTS idx_deductions_student_id ON deductions(student_id);
CREATE INDEX IF NOT EXISTS idx_deductions_stipend_id ON deductions(stipend_id);
CREATE INDEX IF NOT EXISTS idx_deductions_processing_status ON deductions(processing_status);
CREATE INDEX IF NOT EXISTS idx_deductions_deduction_type ON deductions(deduction_type);
CREATE INDEX IF NOT EXISTS idx_deductions_deduction_rule_id ON deductions(deduction_rule_id);

CREATE INDEX IF NOT EXISTS idx_deduction_rules_rule_name ON deduction_rules(rule_name);
CREATE INDEX IF NOT EXISTS idx_deduction_rules_is_active ON deduction_rules(is_active);
CREATE INDEX IF NOT EXISTS idx_deduction_rules_deduction_type ON deduction_rules(deduction_type);

CREATE INDEX IF NOT EXISTS idx_transactions_student_id ON transactions(student_id);
CREATE INDEX IF NOT EXISTS idx_transactions_stipend_id ON transactions(stipend_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_reference_number ON transactions(reference_number);
