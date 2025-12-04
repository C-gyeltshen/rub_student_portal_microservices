-- Fix the banking service schema to match the UUID types
-- Run this before starting the service

-- Drop the existing tables to recreate them with proper schema
DROP TABLE IF EXISTS student_bank_details CASCADE;
DROP TABLE IF EXISTS banks CASCADE;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Recreate banks table with UUID
CREATE TABLE banks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Create index on deleted_at for soft deletes
CREATE INDEX idx_banks_deleted_at ON banks(deleted_at);

-- Recreate student_bank_details table with UUID
CREATE TABLE student_bank_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL,
    bank_id UUID NOT NULL REFERENCES banks(id),
    account_number TEXT NOT NULL UNIQUE,
    account_holder_name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Create index on deleted_at for soft deletes
CREATE INDEX idx_student_bank_details_deleted_at ON student_bank_details(deleted_at);

-- Create index on student_id for faster lookups
CREATE INDEX idx_student_bank_details_student_id ON student_bank_details(student_id);
