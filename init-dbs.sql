-- Initialize separate databases for each microservice
-- This file is automatically executed when PostgreSQL container starts

-- Create database for User Service
CREATE DATABASE user_service_db;

-- Create database for Banking Service
CREATE DATABASE banking_service_db;

-- Create database for Student Management Service
CREATE DATABASE student_service_db;

-- Note: Each service will auto-create its own tables using GORM AutoMigrate