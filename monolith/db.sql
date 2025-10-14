-- ENUM Definitions
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
    hash_password TEXT NOT NULL, -- Removed UNIQUE
    role_id UUID REFERENCES roles(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------------------
-- Student Management Services (Auxiliary)
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
    user_id UUID UNIQUE REFERENCES users(id), -- Added UNIQUE constraint
    name TEXT NOT NULL, 
    rub_id_card_number TEXT NOT NULL UNIQUE, -- Changed INT to TEXT
    email TEXT NOT NULL UNIQUE, -- Corrected INT to TEXT
    phone_number TEXT NOT NULL, -- Corrected INT to TEXT
    date_of_birth DATE, -- Corrected typo and data type
    college_id UUID REFERENCES colleges(id),
    program_id UUID REFERENCES programs(id), -- Added FK to programs
    modified_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stipends(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID REFERENCES students(id),
    amount DECIMAL(10, 2) NOT NULL, -- Changed INT to DECIMAL
    payment_date TIMESTAMPTZ,
    payment_status payment_status_enum, -- Renamed column to avoid type conflict
    payment_method payment_methods_enum,
    journal_number TEXT NOT NULL UNIQUE, -- Changed INT to TEXT/VARCHAR
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------------------
-- Banking Services
-- -----------------------------------------------------------------------

CREATE TABLE banks ( -- Renamed from generic 'accounts' for clarity
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE student_bank_details(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID UNIQUE REFERENCES students(id), -- Added UNIQUE for 1:1 relation
    bank_id UUID REFERENCES banks(ID), -- References 'banks'
    account_number TEXT NOT NULL UNIQUE, -- Corrected INT to TEXT
    account_holder_name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------------------
-- Financial Services
-- -----------------------------------------------------------------------

CREATE TABLE finance_officer (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    officer_user_id UUID UNIQUE REFERENCES users(id), -- Added UNIQUE
    name TEXT NOT NULL, -- Removed UNIQUE
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
    status budget_status_enum,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP 
);

CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    budget_id UUID REFERENCES budget(id),
    requester_id UUID REFERENCES users(id) NOT NULL, -- Assumed requester is a user
    description TEXT NOT NULL, 
    amount DECIMAL(10, 2) NOT NULL, 
    status expenses_status_enum,
    approval_date TIMESTAMPTZ, 
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
); 

CREATE TABLE transaction (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    expenses_id UUID REFERENCES expenses(id), 
    transaction_type transaction_type_enum NOT NULL, 
    amount DECIMAL(12, 2) NOT NULL, 
    transaction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    journal_number TEXT NOT NULL UNIQUE, -- Changed INT to TEXT/VARCHAR
    bank_id UUID REFERENCES banks(id), -- Referenced corrected 'banks' table
    notes TEXT, 
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);