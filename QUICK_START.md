# Quick Start Guide - RUB Student Portal Microservices

## Prerequisites

- Go 1.19 or higher
- PostgreSQL 13+
- Docker & Docker Compose
- Make (optional)
- curl or Postman (for API testing)

## Option 1: Run with Docker Compose (Recommended)

### 1. Clone and Navigate

```bash
cd /path/to/rub_student_portal_microservices
```

### 2. Resolve Git Conflicts (if needed)

```bash
git status
# If there are merge conflicts in docker-compose.yml, resolve them
```

### 3. Start All Services

```bash
docker-compose up -d
```

### 4. Verify Services are Running

```bash
# Check all containers
docker-compose ps

# Expected output:
# rub-postgres-db              ✓ running
# rub-user-services            ✓ running
# rub-banking-services         ✓ running
# rub-student-management-service ✓ running
# rub-finance-services         ✓ running
# rub-api-gateway              ✓ running
```

### 5. Test API Gateway

```bash
curl http://localhost:8080/health
# Response: {"status":"ok"}
```

### 6. Stop Services

```bash
docker-compose down

# Stop and remove all data
docker-compose down -v
```

---

## Option 2: Run Locally (Development)

### Step 1: Start PostgreSQL

```bash
# If using Docker just for PostgreSQL
docker run --name postgres-dev -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 -d postgres:15-alpine

# Create databases
docker exec postgres-dev psql -U postgres -c "
CREATE DATABASE student_service_db;
CREATE DATABASE finance_service_db;
CREATE DATABASE banking_service_db;
CREATE DATABASE user_service_db;
"
```

### Step 2: Update Environment Variables

**Student Management Service** (`services/student_management_service/.env`):

```env
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/student_service_db
PORT=8084
GRPC_PORT=50054
FINANCE_GRPC_URL=localhost:50055
USER_GRPC_URL=localhost:50052
BANKING_GRPC_URL=localhost:50053
```

**Finance Service** (`services/finance_service/.env`):

```env
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/finance_service_db
PORT=8085
GRPC_PORT=50055
STUDENT_GRPC_URL=localhost:50054
```

### Step 3: Start Services in Separate Terminals

**Terminal 1: Student Management Service**

```bash
cd services/student_management_service
go run main.go
```

**Terminal 2: Finance Service**

```bash
cd services/finance_service
go run main.go
```

**Terminal 3: Banking Service**

```bash
cd services/banking_services
go run main.go
```

**Terminal 4: User Service**

```bash
cd services/user_services
go run main.go
```

**Terminal 5: API Gateway**

```bash
cd api-gateway
go run main.go
```

### Step 4: Verify Services

```bash
curl http://localhost:8080/health
```

---

## Service URLs

### API Gateway (Main Entry Point)

- **HTTP**: `http://localhost:8080`
- **Health**: `http://localhost:8080/health`

### Direct Service Access

- **Student Service HTTP**: `http://localhost:8084`
- **Student Service gRPC**: `localhost:50054`
- **Finance Service HTTP**: `http://localhost:8085`
- **Finance Service gRPC**: `localhost:50055`
- **Banking Service HTTP**: `http://localhost:8083`
- **Banking Service gRPC**: `localhost:50053`
- **User Service HTTP**: `http://localhost:8082`
- **User Service gRPC**: `localhost:50052`

### Database

- **PostgreSQL**: `localhost:5432`
- **User**: `postgres`
- **Password**: `postgres` (or `rubpassword` in Docker)

---

## Common Tasks

### Check Service Logs

#### Docker

```bash
# View logs from all services
docker-compose logs

# Follow logs from specific service
docker-compose logs -f student_management_service
docker-compose logs -f finance_services
docker-compose logs -f api-gateway

# View last 100 lines
docker-compose logs --tail=100
```

#### Local Development

```bash
# Logs are printed to stdout in the terminal running the service
```

### Access Database

#### Docker

```bash
docker-compose exec postgres psql -U rubadmin -d student_service_db

# List tables
\dt

# View students
SELECT * FROM students;

# Exit
\q
```

#### Local PostgreSQL

```bash
psql -U postgres -d student_service_db

# Commands same as above
```

### Rebuild Services

#### Docker

```bash
# Rebuild specific service
docker-compose build student_management_service

# Rebuild all services
docker-compose build

# Rebuild and restart
docker-compose up --build -d
```

#### Local Development

```bash
cd services/student_management_service
go build -o student_management_service ./main.go
./student_management_service
```

---

## API Quick Test

### 1. Create Student

```bash
curl -X POST http://localhost:8080/api/students \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "User",
    "student_id": "TEST001",
    "email": "test@example.com",
    "program_id": 1,
    "college_id": 1
  }'
```

### 2. List Students

```bash
curl http://localhost:8080/api/students
```

### 3. Get Specific Student

```bash
curl http://localhost:8080/api/students/1
```

### 4. Calculate Stipend

```bash
curl -X POST http://localhost:8080/api/finance/stipend/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "TEST001",
    "stipend_type": "full-scholarship",
    "amount": 5000
  }'
```

### 5. Create Stipend Record

```bash
curl -X POST http://localhost:8080/api/finance/stipend \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "TEST001",
    "stipend_type": "full-scholarship",
    "amount": 4750,
    "payment_method": "bank_transfer"
  }'
```

### 6. Get Student Stipends

```bash
curl http://localhost:8080/api/finance/stipend/student/TEST001
```

---

## Troubleshooting

### Port Already in Use

```bash
# Find process using port
lsof -i :8080
lsof -i :50054

# Kill process
kill -9 <PID>

# Or use Docker to remove conflicting container
docker ps -a | grep rub
docker rm -f container_name
```

### gRPC Connection Failed

```
Error: failed to connect to finance service: connection refused
```

**Solutions**:

1. Verify Finance Service is running: `docker-compose logs finance_services`
2. Check GRPC_PORT configuration in .env (should be 50055)
3. Ensure services are on same network in Docker
4. For local development, verify port is not blocked by firewall

### Database Connection Failed

```
Error: Failed to connect to database: connection refused
```

**Solutions**:

1. Check PostgreSQL is running: `docker-compose ps postgres`
2. Verify DATABASE_URL in .env files
3. Ensure databases exist: `docker-compose exec postgres psql -U rubadmin -l`
4. Check credentials match

### API Gateway Returns 502

**Solutions**:

1. Verify backend services are running
2. Check logs: `docker-compose logs api-gateway`
3. Ensure service URLs in proxy are correct (should be container names in Docker)

### No Response from Endpoints

```bash
# Check health endpoint first
curl -v http://localhost:8080/health

# Check if routes are registered
# Look for "Routes registered successfully" in logs
docker-compose logs api-gateway | grep "Routes registered"

# Try direct service access
curl http://localhost:8084/api/students
```

---

## Development Workflow

### Making Changes to Student Service

```bash
# 1. Edit code
vim services/student_management_service/handlers/student_handler.go

# 2. Option A: Local restart
#    Press Ctrl+C in running terminal
go run main.go

# 2. Option B: Docker rebuild
docker-compose build student_management_service
docker-compose up -d student_management_service
```

### Adding New Endpoints

1. **Define in Proto** (if using gRPC):

   - Edit `proto/student/student.proto`
   - Add RPC method
   - Generate: `protoc --go_out=. --go-grpc_out=. proto/student/student.proto`

2. **Implement Handler**:

   - Create handler in `handlers/` directory
   - Implement gRPC server method

3. **Register Routes**:

   - Add route in `main.go` or routes file

4. **Test**:
   ```bash
   curl -X POST http://localhost:8084/api/new-endpoint ...
   ```

---

## Data Persistence

### Docker Compose

- All data persists in `postgres_data` volume
- To reset: `docker-compose down -v` (deletes volume)

### Local Development

- PostgreSQL stores data in system PostgreSQL installation
- To reset database:
  ```bash
  psql -U postgres -d student_service_db
  -- Run migration down scripts or delete tables
  ```

---

## Performance Tips

1. **Use indices** on frequently queried columns
2. **Cache expensive calculations** (e.g., stipend calculations)
3. **Batch operations** when processing multiple students
4. **Monitor query performance**:
   ```sql
   -- In PostgreSQL
   EXPLAIN ANALYZE SELECT * FROM students WHERE program_id = 1;
   ```

---

## Production Checklist

- [ ] Update .env with production credentials
- [ ] Enable TLS/mTLS for gRPC communication
- [ ] Add authentication (JWT, OAuth)
- [ ] Configure rate limiting
- [ ] Set up monitoring and logging
- [ ] Configure automated backups
- [ ] Load test all services
- [ ] Set up CI/CD pipeline
- [ ] Enable health checks and auto-restart
- [ ] Document runbooks for incident response

---

## Next Steps

1. Read the full integration guide: `INTEGRATION_GUIDE.md`
2. Review API reference: `API_REFERENCE.md`
3. Check individual service README files
4. Set up monitoring (Prometheus, Grafana)
5. Implement authentication
6. Add comprehensive error handling

---

## Support & Documentation

- **Integration Guide**: See `INTEGRATION_GUIDE.md`
- **API Reference**: See `API_REFERENCE.md`
- **Student Service**: See `services/student_management_service/README.md`
- **Finance Service**: See `services/finance_service/README.md`
- **API Gateway**: See `api-gateway/README.md`

---

## Quick Command Reference

```bash
# Docker Commands
docker-compose up -d                          # Start all services
docker-compose down                           # Stop all services
docker-compose logs -f                        # View all logs
docker-compose ps                             # List services
docker-compose exec postgres psql -U postgres # Access database

# Testing
curl http://localhost:8080/health             # Health check
curl http://localhost:8080/api/students       # Get students

# Development
go run main.go                                 # Run service locally
go test ./...                                  # Run tests
go fmt ./...                                   # Format code
go vet ./...                                   # Check for issues
```

Enjoy! 🚀
