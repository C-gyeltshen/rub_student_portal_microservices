# Finance Service Port Configuration

## Active Services

| Service          | Container            | Host Port | Container Port | Status   | Address                   |
| ---------------- | -------------------- | --------- | -------------- | -------- | ------------------------- |
| PostgreSQL       | rub-postgres         | 5434      | 5432           | Up 4h    | `localhost:5434`          |
| Finance REST API | rub-finance-services | 8084      | 8084           | Up 27m   | `localhost:8084`          |
| gRPC Server      | rub-finance-services | -         | 50051          | Internal | Only in container network |
| API Gateway      | rub-api-gateway      | 8080      | 8080           | Up       | `localhost:8080`          |

## Key Ports

| Port      | Service                   | Type     | Accessibility                       |
| --------- | ------------------------- | -------- | ----------------------------------- |
| **5434**  | PostgreSQL (host machine) | Database | Host & Containers                   |
| **8084**  | Finance Service REST API  | HTTP     | Host & Containers                   |
| **8080**  | API Gateway               | HTTP     | Host & Containers                   |
| **50051** | Finance Service gRPC      | gRPC     | Internal only (not exposed to host) |

## Connection Strings

### From Host Machine (Your Laptop)

```
Database:  postgresql://postgres:postgres@127.0.0.1:5434/rub_student_portal?sslmode=disable
REST API:  http://localhost:8084
API Gate:  http://localhost:8080
```

### From Docker Containers (Internal Network)

```
Database:  postgresql://postgres:postgres@postgres:5432/rub_student_portal?sslmode=disable
gRPC:      localhost:50051
REST API:  http://rub-finance-services:8084
API Gate:  http://rub-api-gateway:8080
```

## Services Overview

### 🐘 PostgreSQL (Port 5434)

- **Container**: rub-postgres
- **Image**: postgres:13
- **Purpose**: Central database for all services
- **Credentials**:
  - User: `postgres`
  - Password: `postgres`
  - Database: `rub_student_portal`

### 💰 Finance Service (Port 8084)

- **Container**: rub-finance-services
- **Purpose**: Handle stipends, deductions, and financial calculations
- **Endpoints**:
  - REST API: `http://localhost:8084`
  - gRPC: `localhost:50051` (internal)
- **Features**:
  - Stipend Management
  - Deduction Rules
  - Financial Calculations

### 🌐 API Gateway (Port 8080)

- **Container**: rub-api-gateway
- **Purpose**: Central entry point for all microservices
- **Base URL**: `http://localhost:8080`

## Environment Variables for Tests

```bash
# Database connection for host machine
export DATABASE_URL="postgresql://postgres:postgres@127.0.0.1:5434/rub_student_portal?sslmode=disable"

# Database connection for containers
export DATABASE_URL="postgresql://postgres:postgres@postgres:5432/rub_student_portal?sslmode=disable"

# gRPC port (internal)
export GRPC_PORT="50051"

# HTTP port
export PORT="8084"
```

## Docker Compose Configuration

```yaml
postgres:
  ports:
    - "5434:5432" # Maps host:5434 to container:5432

finance_services:
  ports:
    - "8084:8084" # Maps host:8084 to container:8084
    # gRPC port 50051 NOT exposed (internal only)

api-gateway:
  ports:
    - "8080:8080" # Maps host:8080 to container:8080
```

## Testing Port Connectivity

### Check PostgreSQL Connection

```bash
psql -h 127.0.0.1 -p 5434 -U postgres -d rub_student_portal
```

### Check Finance Service REST API

```bash
curl http://localhost:8084/health
```

### Check API Gateway

```bash
curl http://localhost:8080/health
```

### Check gRPC Service (from container)

```bash
# Must be run from inside a container on rub-network
grpcurl -plaintext localhost:50051 list
```

## Network Configuration

All services are connected via the `rub-network` bridge network, allowing them to communicate using container names as hostnames.

### DNS Resolution:

- `postgres:5432` - PostgreSQL database
- `rub-finance-services:8084` - Finance service
- `rub-api-gateway:8080` - API gateway

## Important Notes

⚠️ **gRPC Port (50051)** is **NOT** exposed to the host machine. It only exists within the Docker container network for internal service-to-service communication.

✅ **REST API (8084)** is fully accessible from both the host machine and other containers.

✅ **Database (5434)** is accessible from the host machine at port 5434, but inside containers it connects to the internal port 5432 using the hostname `postgres`.
