package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"api-gateway/internal/config"
)

// UserContext represents the authenticated user information
type UserContext struct {
	UID           string            `json:"uid"`
	Email         string            `json:"email"`
	EmailVerified bool              `json:"email_verified"`
	Role          string            `json:"role"`
	CollegeID     string            `json:"college_id,omitempty"`
	Permissions   []string          `json:"permissions,omitempty"`
	CustomClaims  map[string]interface{} `json:"custom_claims,omitempty"`
}

// ContextKey is the type used for context keys
type ContextKey string

const (
	// UserContextKey is the context key for user information
	UserContextKey ContextKey = "user"
)

// Role constants
const (
	RoleAdmin          = "admin"
	RoleFinanceOfficer = "finance_officer"
	RoleStudent        = "student"
)

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Success   bool      `json:"success"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Code      string    `json:"code"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path,omitempty"`
}

// WriteErrorResponse writes a standardized error response
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message, code string) {
	response := ErrorResponse{
		Success:   false,
		Error:     http.StatusText(statusCode),
		Message:   message,
		Code:      code,
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to write error response: %v", err)
	}
}

// AuthMiddleware verifies Firebase ID tokens and attaches user information to context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("Missing Authorization header for %s %s", r.Method, r.URL.Path)
			WriteErrorResponse(w, http.StatusUnauthorized, "Missing Authorization header", "AUTH_001")
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("Invalid Authorization header format for %s %s", r.Method, r.URL.Path)
			WriteErrorResponse(w, http.StatusUnauthorized, "Invalid Authorization header format", "AUTH_002")
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			log.Printf("Empty token in Authorization header for %s %s", r.Method, r.URL.Path)
			WriteErrorResponse(w, http.StatusUnauthorized, "Empty token", "AUTH_003")
			return
		}

		// Verify token with Firebase
		authClient := config.GetAuthClient()
		if authClient == nil {
			log.Printf("Firebase Auth client not initialized")
			WriteErrorResponse(w, http.StatusInternalServerError, "Authentication service unavailable", "AUTH_004")
			return
		}

		firebaseToken, err := authClient.VerifyIDToken(context.Background(), token)
		if err != nil {
			log.Printf("Token verification failed for %s %s: %v", r.Method, r.URL.Path, err)
			if strings.Contains(err.Error(), "expired") {
				WriteErrorResponse(w, http.StatusUnauthorized, "Token expired", "AUTH_005")
			} else if strings.Contains(err.Error(), "invalid") {
				WriteErrorResponse(w, http.StatusForbidden, "Invalid token", "AUTH_006")
			} else {
				WriteErrorResponse(w, http.StatusForbidden, "Token verification failed", "AUTH_007")
			}
			return
		}

		// Extract user information
		userCtx := &UserContext{
			UID:           firebaseToken.UID,
			Email:         "",
			EmailVerified: false,
			CustomClaims:  firebaseToken.Claims,
		}

		// Get additional user info from Firebase
		userRecord, err := authClient.GetUser(context.Background(), firebaseToken.UID)
		if err != nil {
			log.Printf("Failed to get user record for UID %s: %v", firebaseToken.UID, err)
			// Continue with limited user info
		} else {
			userCtx.Email = userRecord.Email
			userCtx.EmailVerified = userRecord.EmailVerified
		}

		// Extract custom claims
		if role, ok := firebaseToken.Claims["role"].(string); ok {
			userCtx.Role = role
		}

		if collegeID, ok := firebaseToken.Claims["college_id"].(string); ok {
			userCtx.CollegeID = collegeID
		}

		if permissions, ok := firebaseToken.Claims["permissions"].([]interface{}); ok {
			userCtx.Permissions = make([]string, len(permissions))
			for i, perm := range permissions {
				if permStr, ok := perm.(string); ok {
					userCtx.Permissions[i] = permStr
				}
			}
		}

		// Log successful authentication
		log.Printf("User authenticated successfully: UID=%s, Email=%s, Role=%s", userCtx.UID, userCtx.Email, userCtx.Role)

		// Add user context to request
		ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
		r = r.WithContext(ctx)

		// Call next handler
		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext extracts user information from request context
func GetUserFromContext(ctx context.Context) (*UserContext, bool) {
	user, ok := ctx.Value(UserContextKey).(*UserContext)
	return user, ok
}

// AuthorizeRoles creates a middleware that checks if user has any of the required roles
func AuthorizeRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context (should be set by AuthMiddleware)
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				log.Printf("No user context found for authorization check on %s %s", r.Method, r.URL.Path)
				WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "AUTHZ_001")
				return
			}

			// Check if user has required role
			userRole := user.Role
			if userRole == "" {
				log.Printf("User %s has no role assigned for %s %s", user.UID, r.Method, r.URL.Path)
				WriteErrorResponse(w, http.StatusForbidden, "No role assigned", "AUTHZ_002")
				return
			}

			// Check if user's role is in allowed roles
			hasRequiredRole := false
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					hasRequiredRole = true
					break
				}
			}

			if !hasRequiredRole {
				log.Printf("User %s with role %s attempted to access %s %s requiring roles: %v", 
					user.UID, userRole, r.Method, r.URL.Path, allowedRoles)
				WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions", "AUTHZ_003")
				return
			}

			log.Printf("Authorization successful: User %s (role: %s) accessing %s %s", 
				user.UID, userRole, r.Method, r.URL.Path)

			next.ServeHTTP(w, r)
		})
	}
}

// AuthorizeAdmin creates a middleware that only allows admin users
func AuthorizeAdmin() func(http.Handler) http.Handler {
	return AuthorizeRoles(RoleAdmin)
}

// AuthorizeFinanceOfficer creates a middleware that allows admin and finance officer users
func AuthorizeFinanceOfficer() func(http.Handler) http.Handler {
	return AuthorizeRoles(RoleAdmin, RoleFinanceOfficer)
}

// AuthorizeStudent creates a middleware that allows all authenticated users (admin, finance_officer, student)
func AuthorizeStudent() func(http.Handler) http.Handler {
	return AuthorizeRoles(RoleAdmin, RoleFinanceOfficer, RoleStudent)
}