# RUB Student Portal Microservices Architecture

## Overview

This repository contains a microservices-based implementation of the Royal University of Bhutan Student Portal system. The architecture employs a distributed system design pattern, decomposing the monolithic application into independently deployable services that communicate through well-defined APIs and gRPC protocols.

## System Architecture

The system follows a microservices architecture pattern with the following core components:

### API Gateway
The API Gateway serves as the single entry point for all client requests, implementing cross-cutting concerns such as routing, authentication, and request aggregation. It proxies requests to the appropriate backend microservices and handles CORS policies.

### Core Microservices

#### Student Management Service
Manages student registration, enrollment, academic records, and profile information. This service implements comprehensive CRUD operations for student entities and integrates with the finance service for fee-related operations. The service supports bulk operations for efficient data management through CSV imports.

#### Banking Services
Handles banking and financial transaction data, including bank information management and student bank account details. The service maintains relationships between students and their associated banking information required for financial operations.

#### Finance Service
Manages all financial aspects of the student lifecycle, including fee structures, payment processing, financial aid, scholarships, and revenue tracking. The service implements a comprehensive audit trail system and provides detailed financial reporting capabilities.

#### User Services
Provides authentication, authorization, and user management functionality. This service handles user credentials, role-based access control, and session management across the platform.

## Technical Stack

### Backend Technologies
- **Primary Language**: Go (Golang)
- **Communication Protocols**: REST API, gRPC
- **Database**: PostgreSQL (relational database for data persistence)
- **Containerization**: Docker, Docker Compose

### Architecture Patterns
- Microservices Architecture
- API Gateway Pattern
- Database Per Service Pattern
- Protocol Buffers for service contracts
- RESTful API design principles

## Project Structure

The repository is organized into the following directories:

- `api-gateway/`: API Gateway implementation with routing and middleware
- `services/`: Individual microservice implementations
  - `student_management_service/`: Student data and enrollment management
  - `banking_services/`: Banking and account information management
  - `finance_service/`: Financial operations and fee management
  - `user_services/`: Authentication and user management
- `proto/`: Protocol Buffer definitions for gRPC communication
- `monolith/`: Legacy monolithic implementation (for reference)

## Development Setup

### Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL 14 or higher
- Make utility

### Building and Running

The project includes a Makefile for common development tasks and a docker-compose configuration for orchestrating the microservices environment.

Refer to `QUICK_START.md` for detailed setup instructions and `CSV_UPLOAD_GUIDE.md` for bulk data import procedures.

## Testing

Each microservice includes unit tests and integration tests. Test coverage reports are generated during the build process. Refer to `TESTING_SUMMARY.md` for comprehensive testing documentation.

## API Documentation

Individual services contain their own API documentation:
- Finance Service: See `services/finance_service/API_REFERENCE.md` and `ENDPOINTS_AND_CRUD_MAPPING.md`
- Student Management Service: See `services/student_management_service/POSTMAN_TESTING_GUIDE.md`

## Database Management

Database schemas are defined per service:
- Finance Service: `services/finance_service/db_schema.sql`
- Student Management Service: `services/student_management_service/db_schema.sql`
- Banking Services: `services/banking_services/fix_schema.sql`

Initial database setup scripts are provided in `init-dbs.sql` and `db.sql`.

## Integration Points

The microservices communicate through:
1. Synchronous REST API calls for client-facing operations
2. gRPC for inter-service communication requiring low latency
3. Database transactions within service boundaries
4. Shared protocol definitions in the `proto/` directory

## Service Ports Configuration

Port allocation for individual services is documented in `services/finance_service/PORTS_CONFIGURATION.md`.

## Academic Context

This project serves as a practical implementation of modern software engineering principles, demonstrating:
- Microservices architecture design and implementation
- Service decomposition strategies
- Inter-service communication patterns
- Database design and management in distributed systems
- Containerization and orchestration
- API design and documentation
- Testing strategies for distributed systems

## License

This project is developed for academic purposes at the Royal University of Bhutan.