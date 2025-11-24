-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE payment_status_enum AS ENUM ('Pending', 'Processed', 'Failed');
CREATE TYPE payment_methods_enum AS ENUM ('Bank_transfer', 'E-payment');
CREATE TYPE master_status_enum AS ENUM ('Enabled', 'Disabled');
CREATE TYPE question_type_enum AS ENUM('Radio', 'Checkbox');

-- Financial Services Enums
CREATE TYPE budget_status_enum AS ENUM('Active', 'Close', 'Expire');
CREATE TYPE expenses_status_enum AS ENUM('Pending', 'Approved', 'Rejected', 'Paid');
CREATE TYPE transaction_type_enum AS ENUM('DEBIT', 'CREDIT');

-- -----------------------------------------------------------------------
-- User Services
-- -----------------------------------------------------------------------

CREATE TABLE roles(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    hash_password TEXT NOT NULL,
    role_id UUID REFERENCES roles(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------------------
-- Student Management Services
-- -----------------------------------------------------------------------

CREATE TABLE colleges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    modified_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE programs(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    college_id UUID REFERENCES colleges(id),
    modified_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE students(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE REFERENCES users(id),
    name TEXT NOT NULL, 
    rub_id_card_number TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    phone_number TEXT NOT NULL,
    date_of_birth DATE,
    college_id UUID REFERENCES colleges(id),
    program_id UUID REFERENCES programs(id),
    modified_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------------------
-- Banking Services
-- -----------------------------------------------------------------------

CREATE TABLE banks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE student_bank_details(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID UNIQUE REFERENCES students(id),
    bank_id UUID REFERENCES banks(id),
    account_number TEXT NOT NULL UNIQUE,
    account_holder_name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------------------
-- Finance Service - Stipend Calculation & Deduction System
-- -----------------------------------------------------------------------

-- Stipends Table
CREATE TABLE stipends(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount >= 0),
    stipend_type VARCHAR(50) NOT NULL,
    payment_date TIMESTAMPTZ,
    payment_status VARCHAR(50) DEFAULT 'Pending',
    payment_method VARCHAR(50),
    journal_number TEXT NOT NULL UNIQUE,
    transaction_id UUID,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Deduction Rules Table
CREATE TABLE deduction_rules(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rule_name VARCHAR(100) NOT NULL UNIQUE,
    deduction_type VARCHAR(100) NOT NULL,
    description TEXT,
    base_amount DECIMAL(10, 2) NOT NULL CHECK (base_amount >= 0),
    max_deduction_amount DECIMAL(10, 2) NOT NULL CHECK (max_deduction_amount >= 0),
    min_deduction_amount DECIMAL(10, 2) DEFAULT 0 CHECK (min_deduction_amount >= 0),
    is_applicable_to_full_scholar BOOLEAN DEFAULT FALSE,
    is_applicable_to_self_funded BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    applies_monthly BOOLEAN DEFAULT FALSE,
    applies_annually BOOLEAN DEFAULT FALSE,
    is_optional BOOLEAN DEFAULT FALSE,
    priority INTEGER DEFAULT 0,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_by UUID REFERENCES users(id) ON DELETE SET NULL,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Deductions Table
CREATE TABLE deductions(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    deduction_rule_id UUID NOT NULL REFERENCES deduction_rules(id) ON DELETE RESTRICT,
    stipend_id UUID NOT NULL REFERENCES stipends(id) ON DELETE CASCADE,
    amount DECIMAL(10, 2) NOT NULL CHECK (amount >= 0),
    deduction_type VARCHAR(100) NOT NULL,
    description TEXT,
    deduction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    processing_status VARCHAR(50) DEFAULT 'Pending',
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approval_date TIMESTAMPTZ,
    rejection_reason TEXT,
    transaction_id UUID,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------------------
-- Financial Services (Budget & Expenses)
-- -----------------------------------------------------------------------

CREATE TABLE finance_officer (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    officer_user_id UUID UNIQUE REFERENCES users(id),
    name TEXT NOT NULL,
    college_id UUID REFERENCES colleges(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE budget (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    finance_officer_id UUID REFERENCES finance_officer(id),
    name TEXT NOT NULL, 
    purpose TEXT,
    allocated_amount DECIMAL(12, 2) NOT NULL, 
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP 
);

CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    budget_id UUID REFERENCES budget(id),
    requester_id UUID REFERENCES users(id) NOT NULL,
    description TEXT NOT NULL, 
    amount DECIMAL(10, 2) NOT NULL, 
    status VARCHAR(50),
    approval_date TIMESTAMPTZ, 
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
); 

CREATE TABLE transaction (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    expenses_id UUID REFERENCES expenses(id), 
    transaction_type VARCHAR(50) NOT NULL, 
    amount DECIMAL(12, 2) NOT NULL, 
    transaction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    journal_number TEXT NOT NULL UNIQUE,
    bank_id UUID REFERENCES banks(id),
    notes TEXT, 
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------------------
-- Indexes for Performance
-- -----------------------------------------------------------------------

CREATE INDEX idx_stipends_student_id ON stipends(student_id);
CREATE INDEX idx_stipends_payment_status ON stipends(payment_status);
CREATE INDEX idx_stipends_stipend_type ON stipends(stipend_type);
CREATE INDEX idx_stipends_created_at ON stipends(created_at);

CREATE INDEX idx_deductions_student_id ON deductions(student_id);
CREATE INDEX idx_deductions_stipend_id ON deductions(stipend_id);
CREATE INDEX idx_deductions_processing_status ON deductions(processing_status);
CREATE INDEX idx_deductions_deduction_type ON deductions(deduction_type);
CREATE INDEX idx_deductions_deduction_rule_id ON deductions(deduction_rule_id);

CREATE INDEX idx_deduction_rules_rule_name ON deduction_rules(rule_name);
CREATE INDEX idx_deduction_rules_is_active ON deduction_rules(is_active);
CREATE INDEX idx_deduction_rules_deduction_type ON deduction_rules(deduction_type);



