# Database Setup Guide

## Architecture Overview

Each microservice has its **own separate database** on a shared PostgreSQL server:

```
PostgreSQL Server (port 5432)
├── user_service_db       → User Service (port 8082)
├── banking_service_db    → Banking Service (port 8083)
└── student_service_db    → Student Management Service (port 8084)
```

## How to Run

### Option 1: Using Docker (Recommended - Easiest!)

```bash
# Start all services with database
docker-compose up -d

# Check if services are running
docker-compose ps

# View logs
docker-compose logs -f student_management_service

# Stop all services
docker-compose down
```

That's it! Docker will automatically:

- Create PostgreSQL server
- Create 3 separate databases (user_service_db, banking_service_db, student_service_db)
- Start all 3 microservices
- Each service connects to its own database

### Option 2: Local Development (Without Docker)

**Step 1: Install PostgreSQL**

```bash
brew install postgresql
brew services start postgresql
```

**Step 2: Create the databases**

```bash
# Connect to PostgreSQL
psql postgres

# Run these commands in psql:
CREATE USER rubadmin WITH PASSWORD 'rubpassword';
CREATE DATABASE user_service_db OWNER rubadmin;
CREATE DATABASE banking_service_db OWNER rubadmin;
CREATE DATABASE student_service_db OWNER rubadmin;
\q
```

**Step 3: Run each service**

```bash
# Terminal 1: User Service
cd services/user_services
go run main.go

# Terminal 2: Banking Service
cd services/banking_services
go run main.go

# Terminal 3: Student Management Service
cd services/student_management_service
go run main.go

# Terminal 4: API Gateway
cd api-gateway
go run main.go
```

## Database Credentials

- **Username**: `rubadmin`
- **Password**: `rubpassword`
- **Host**: `localhost` (local) or `postgres` (Docker)
- **Port**: `5432`

### Connection Strings

**Student Service:**

- Local: `postgresql://rubadmin:rubpassword@localhost:5432/student_service_db`
- Docker: `postgresql://rubadmin:rubpassword@postgres:5432/student_service_db`

**User Service:**

- Local: `postgresql://rubadmin:rubpassword@localhost:5432/user_service_db`
- Docker: `postgresql://rubadmin:rubpassword@postgres:5432/user_service_db`

**Banking Service:**

- Local: `postgresql://rubadmin:rubpassword@localhost:5432/banking_service_db`
- Docker: `postgresql://rubadmin:rubpassword@postgres:5432/banking_service_db`

## How Tables Are Created

**No manual SQL needed!** Each service uses GORM AutoMigrate to automatically create tables:

### Student Service Creates:

- `students`
- `colleges`
- `programs`
- `stipend_allocations`
- `stipend_histories`
- `audit_logs`

### User Service Creates:

- `users`
- `roles`

### Banking Service Creates:

- `banks`
- `student_bank_details`

## Testing the Setup

```bash
# Check if database is running
docker-compose ps postgres

# Connect to database
docker exec -it rub-postgres-db psql -U rubadmin -d student_service_db

# List all tables
\dt

# Exit psql
\q

# Test API
curl http://localhost:8084/api/students
```

## Troubleshooting

**Problem: "DATABASE_URL not set"**

- Make sure `.env` file exists in each service folder
- Check that DATABASE_URL is set correctly

**Problem: "Cannot connect to database"**

- Check if PostgreSQL is running: `docker-compose ps`
- Verify credentials in `.env` file

**Problem: "No tables created"**

- Tables are auto-created when service starts
- Check service logs: `docker-compose logs student_management_service`

**Problem: Port already in use**

```bash
# Find what's using the port
lsof -i :5432

# Kill the process
kill -9 <PID>
```

## Project Structure

```
rub_student_portal_microservices/
├── docker-compose.yml          # Orchestrates all services + database
├── init-dbs.sql               # Creates 3 databases on startup
├── services/
│   ├── user_services/
│   │   ├── .env              # DATABASE_URL=...user_service_db
│   │   └── database/db.go    # Connects to user_service_db
│   ├── banking_services/
│   │   ├── .env              # DATABASE_URL=...banking_service_db
│   │   └── database/db.go    # Connects to banking_service_db
│   └── student_management_service/
│       ├── .env              # DATABASE_URL=...student_service_db
│       └── database/db.go    # Connects to student_service_db
└── api-gateway/
```

## Next Steps

1. **Start the services**: `docker-compose up -d`
2. **Test Student API**: `curl http://localhost:8084/api/students`
3. **Create a student**: Use API_DOCUMENTATION.md for examples
4. **View data**: Use pgAdmin or `psql` to inspect tables
