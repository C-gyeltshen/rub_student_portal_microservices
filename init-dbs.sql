-- Initialize separate databases for each microservice
-- This file is automatically executed when PostgreSQL container starts

-- Create database for User Service
CREATE DATABASE user_service_db;

-- Create database for Banking Service
CREATE DATABASE banking_service_db;

-- Create database for Student Management Service
CREATE DATABASE student_service_db;

-- Create database for Finance Service
CREATE DATABASE finance_service_db;

-- Note: Each service will auto-create its own tables using GORM AutoMigrate

-- Finance Service - Audit Logs Table (manual creation for reliability)
\c finance_service_db
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    finance_officer VARCHAR(255),
    description TEXT,
    old_values JSONB DEFAULT '{}',
    new_values JSONB DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'SUCCESS',
    error_message TEXT,
    ip_address VARCHAR(50),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_finance_officer ON audit_logs(finance_officer);