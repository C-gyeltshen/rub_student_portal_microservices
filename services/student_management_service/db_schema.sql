-- Student Management Service Database Schema
-- This schema matches the GORM models defined in the models package
-- Uses auto-incrementing integer IDs and GORM conventions

-- Drop tables if they exist (for fresh installations)
DROP TABLE IF EXISTS stipend_histories CASCADE;
DROP TABLE IF EXISTS stipend_allocations CASCADE;
DROP TABLE IF EXISTS students CASCADE;
DROP TABLE IF EXISTS programs CASCADE;
DROP TABLE IF EXISTS colleges CASCADE;

-- -----------------------------------------------------------------------
-- Colleges Table
-- -----------------------------------------------------------------------
CREATE TABLE colleges (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create index on deleted_at for soft delete queries
CREATE INDEX idx_colleges_deleted_at ON colleges(deleted_at);

-- -----------------------------------------------------------------------
-- Programs Table
-- -----------------------------------------------------------------------
CREATE TABLE programs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    college_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (college_id) REFERENCES colleges(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_programs_deleted_at ON programs(deleted_at);
CREATE INDEX idx_programs_college_id ON programs(college_id);

-- -----------------------------------------------------------------------
-- Students Table
-- -----------------------------------------------------------------------
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    rub_id_card_number VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(255),
    date_of_birth VARCHAR(255),
    program_id INTEGER,
    college_id INTEGER,
    user_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (program_id) REFERENCES programs(id) ON DELETE SET NULL,
    FOREIGN KEY (college_id) REFERENCES colleges(id) ON DELETE SET NULL
);

-- Create indexes
CREATE INDEX idx_students_deleted_at ON students(deleted_at);
CREATE INDEX idx_students_program_id ON students(program_id);
CREATE INDEX idx_students_college_id ON students(college_id);
CREATE INDEX idx_students_user_id ON students(user_id);
CREATE INDEX idx_students_rub_id_card_number ON students(rub_id_card_number);
CREATE INDEX idx_students_email ON students(email);

-- -----------------------------------------------------------------------
-- Stipend Allocations Table
-- Finance Service Integration - Stores stipend allocation records
-- -----------------------------------------------------------------------
CREATE TABLE stipend_allocations (
    id SERIAL PRIMARY KEY,
    allocation_id VARCHAR(255) UNIQUE NOT NULL,
    student_id INTEGER NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    allocation_date VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending',
    approved_by INTEGER,
    approval_date VARCHAR(255),
    semester INTEGER,
    academic_year VARCHAR(255),
    remarks TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
);

-- Create indexes for stipend_allocations
CREATE INDEX idx_stipend_allocations_student_id ON stipend_allocations(student_id);
CREATE INDEX idx_stipend_allocations_status ON stipend_allocations(status);
CREATE INDEX idx_stipend_allocations_allocation_id ON stipend_allocations(allocation_id);
CREATE INDEX idx_stipend_allocations_deleted_at ON stipend_allocations(deleted_at);

-- -----------------------------------------------------------------------
-- Stipend Histories Table
-- Finance Service Integration - Stores stipend payment history
-- -----------------------------------------------------------------------
CREATE TABLE stipend_histories (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(255) UNIQUE NOT NULL,
    student_id INTEGER NOT NULL,
    allocation_id INTEGER,
    amount DECIMAL(10,2) NOT NULL,
    payment_date VARCHAR(255),
    transaction_status VARCHAR(50),
    payment_method VARCHAR(50),
    bank_reference VARCHAR(255),
    remarks TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    FOREIGN KEY (allocation_id) REFERENCES stipend_allocations(id) ON DELETE SET NULL
);

-- Create indexes for stipend_histories
CREATE INDEX idx_stipend_histories_student_id ON stipend_histories(student_id);
CREATE INDEX idx_stipend_histories_allocation_id ON stipend_histories(allocation_id);
CREATE INDEX idx_stipend_histories_transaction_id ON stipend_histories(transaction_id);
CREATE INDEX idx_stipend_histories_transaction_status ON stipend_histories(transaction_status);
CREATE INDEX idx_stipend_histories_deleted_at ON stipend_histories(deleted_at);

-- -----------------------------------------------------------------------
-- Comments
-- -----------------------------------------------------------------------
COMMENT ON TABLE colleges IS 'Stores academic colleges/institutions';
COMMENT ON TABLE programs IS 'Stores academic programs offered by colleges';
COMMENT ON TABLE students IS 'Stores student information and their academic affiliations';
COMMENT ON TABLE stipend_allocations IS 'Finance service integration - tracks stipend allocations for students';
COMMENT ON TABLE stipend_histories IS 'Finance service integration - tracks stipend payment history';

COMMENT ON COLUMN students.name IS 'Student full name';
COMMENT ON COLUMN students.rub_id_card_number IS 'Unique RUB student ID card number';
COMMENT ON COLUMN students.user_id IS 'Reference to user in user service (not enforced by FK)';
