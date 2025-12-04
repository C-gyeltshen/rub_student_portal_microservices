-- Fix the students table schema
-- Run this SQL script directly in your cloud database to fix the schema issue

-- Drop dependent tables first (in correct order)
DROP TABLE IF EXISTS stipend_histories CASCADE;
DROP TABLE IF EXISTS stipend_allocations CASCADE;
DROP TABLE IF EXISTS students CASCADE;
DROP TABLE IF EXISTS programs CASCADE;
DROP TABLE IF EXISTS colleges CASCADE;

-- Recreate colleges table
CREATE TABLE colleges (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_colleges_deleted_at ON colleges(deleted_at);

-- Recreate programs table
CREATE TABLE programs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    college_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (college_id) REFERENCES colleges(id) ON DELETE CASCADE
);

CREATE INDEX idx_programs_deleted_at ON programs(deleted_at);
CREATE INDEX idx_programs_college_id ON programs(college_id);

-- Recreate students table (WITHOUT student_id column, only id)
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

CREATE INDEX idx_students_deleted_at ON students(deleted_at);
CREATE INDEX idx_students_program_id ON students(program_id);
CREATE INDEX idx_students_college_id ON students(college_id);
CREATE INDEX idx_students_user_id ON students(user_id);
CREATE INDEX idx_students_rub_id_card_number ON students(rub_id_card_number);
CREATE INDEX idx_students_email ON students(email);

-- Recreate stipend_allocations table
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

CREATE INDEX idx_stipend_allocations_student_id ON stipend_allocations(student_id);
CREATE INDEX idx_stipend_allocations_status ON stipend_allocations(status);
CREATE INDEX idx_stipend_allocations_allocation_id ON stipend_allocations(allocation_id);
CREATE INDEX idx_stipend_allocations_deleted_at ON stipend_allocations(deleted_at);

-- Recreate stipend_histories table
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

CREATE INDEX idx_stipend_histories_student_id ON stipend_histories(student_id);
CREATE INDEX idx_stipend_histories_allocation_id ON stipend_histories(allocation_id);
CREATE INDEX idx_stipend_histories_transaction_id ON stipend_histories(transaction_id);
CREATE INDEX idx_stipend_histories_transaction_status ON stipend_histories(transaction_status);
CREATE INDEX idx_stipend_histories_deleted_at ON stipend_histories(deleted_at);
