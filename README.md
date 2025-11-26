# RUB Student Portal Microservices

A microservices-based student portal system for Royal University of Bhutan with Firebase authentication and authorization.

## Architecture Overview

This system uses a microservices architecture with the following components:

### Services:

- **API Gateway** (Port 8080): Single entry point with Firebase authentication
- **User Services** (Port 8082): User and role management
- **Banking Services** (Port 8083): Bank and student banking details

### Authentication & Authorization:

- **Firebase Authentication**: Token-based authentication using Firebase ID tokens
- **Role-Based Access Control**: Three roles (admin, finance_officer, student)
- **Custom Claims**: Firebase custom claims for permissions and college affiliation
- **Route Protection**: All API routes protected except public endpoints

## Features

- ✅ Firebase Authentication integration
- ✅ Role-based authorization (RBAC)
- ✅ Secure token verification
- ✅ User context forwarding to downstream services
- ✅ Comprehensive error handling and security logging
- ✅ Custom claims management
- ✅ College-based access control

## Quick Start

### 1. Setup Firebase

Follow the detailed instructions in [FIREBASE_SETUP.md](FIREBASE_SETUP.md)

### 2. Configure Environment

```bash
cd api-gateway
cp .env.example .env
# Edit .env with your Firebase credentials
```

### 3. Start Services

```bash
docker-compose up --build
```

### 4. Test Authentication

```bash
# Test public endpoint
curl http://localhost:8080/health

# Test protected endpoint (requires Firebase token)
curl -H "Authorization: Bearer YOUR_FIREBASE_ID_TOKEN" \
     http://localhost:8080/profile
```

## API Endpoints

### Public Endpoints (No Authentication Required)

- `GET /` - Login instructions
- `GET /dashboard` - Public dashboard
- `GET /health` - Health check
- `GET /login-instructions` - Authentication instructions

### Protected Endpoints (Require Authentication)

- `GET /profile` - User profile (all authenticated users)
- `GET /admin/dashboard` - Admin dashboard (admin only)
- `GET /finance/dashboard` - Finance dashboard (admin, finance_officer)
- `GET /student/dashboard` - Student dashboard (all users)

### API Routes (Role-Based Access)

- `GET|POST|PATCH|DELETE /api/users/*` - User management (admin, finance_officer)
- `GET /api/banks/*` - View banks (all users)
- `POST|PUT|DELETE /api/banks/*` - Manage banks (admin, finance_officer)
- `GET|POST|PUT /api/student-bank-details/*` - Student banking (role-dependent)

## User Roles & Permissions

### Admin (`admin`)

- Full system access
- All CRUD operations
- User and role management
- Cross-college access

### Finance Officer (`finance_officer`)

- Budget and expense management
- Student banking oversight
- College-specific access
- User viewing permissions

### Student (`student`)

- View own profile and data
- Manage own banking details
- View available banks
- College-specific access

## Security Features

- **Token Verification**: Firebase ID tokens verified on every request
- **Custom Claims**: Role and permission enforcement
- **Request Logging**: Comprehensive security event logging
- **Error Handling**: Structured error responses with codes
- **College Isolation**: Users can only access their college data
- **Header Forwarding**: User context passed to downstream services

## Development

### Running Locally

```bash
# Start individual services
cd user_services && go run main.go
cd banking_services && go run main.go
cd api-gateway && go run main.go
```

### Environment Variables

Required Firebase configuration:

```env
FIREBASE_PROJECT_ID=your-project-id
FIREBASE_CLIENT_EMAIL=your-service-account-email
FIREBASE_PRIVATE_KEY=your-private-key
```

### Testing

```bash
# Run security scan (if Snyk is available)
snyk code test

# Test with different roles
curl -H "Authorization: Bearer ADMIN_TOKEN" http://localhost:8080/admin/dashboard
curl -H "Authorization: Bearer STUDENT_TOKEN" http://localhost:8080/student/dashboard
```

## Security Considerations

1. **Never expose service account credentials** in client code
2. **Use HTTPS** in production environments
3. **Rotate Firebase service account keys** regularly
4. **Monitor authentication logs** for suspicious activity
5. **Implement token refresh** in frontend applications

## Contributing

1. Fork the repository
2. Create a feature branch
3. Follow security best practices
4. Test authentication and authorization
5. Submit a pull request

## License

This project is licensed under the MIT License.
