package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

// RegisterRequest represents the registration payload
type RegisterRequest struct {
	UID  string `json:"uid"`
	Role string `json:"role"`
}

// SetRoleRequest represents the role assignment payload
type SetRoleRequest struct {
	UID  string `json:"uid"`
	Role string `json:"role"`
}

// Valid roles
var ValidRoles = []string{"admin", "finance_officer", "student"}

// RegisterUserHandler handles user registration with role assignment
func RegisterUserHandler(authClient *auth.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Invalid request body for registration: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		
		// Validate required fields
		if req.UID == "" {
			http.Error(w, "UID is required", http.StatusBadRequest)
			return
		}
		
		// Default to student role if not specified
		if req.Role == "" {
			req.Role = "student"
		}
		
		// Validate role
		if !isValidRole(req.Role) {
			http.Error(w, "Invalid role. Allowed roles: student", http.StatusBadRequest)
			return
		}
		
		// Security: Only allow student registration for self-service
		if req.Role != "student" {
			http.Error(w, "Only student role allowed for self-registration", http.StatusForbidden)
			return
		}
		
		// Verify user exists in Firebase
		_, err := authClient.GetUser(context.Background(), req.UID)
		if err != nil {
			log.Printf("Error verifying user %s: %v", req.UID, err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		
		// Set custom claims
		claims := map[string]interface{}{
			"role": req.Role,
		}
		
		err = authClient.SetCustomUserClaims(context.Background(), req.UID, claims)
		if err != nil {
			log.Printf("Error setting custom claims for user %s: %v", req.UID, err)
			http.Error(w, "Failed to assign role", http.StatusInternalServerError)
			return
		}
		
		log.Printf("Role '%s' assigned to user %s during registration", req.Role, req.UID)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User registered successfully",
			"uid":     req.UID,
			"role":    req.Role,
		})
	}
}

// SetUserRoleHandler sets custom claims for a user (admin only)
func SetUserRoleHandler(authClient *auth.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SetRoleRequest
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Invalid request body for role assignment: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		
		// Validate required fields
		if req.UID == "" || req.Role == "" {
			http.Error(w, "UID and role are required", http.StatusBadRequest)
			return
		}
		
		// Validate role
		if !isValidRole(req.Role) {
			http.Error(w, "Invalid role. Allowed roles: admin, finance_officer, student", http.StatusBadRequest)
			return
		}
		
		// Verify user exists in Firebase
		user, err := authClient.GetUser(context.Background(), req.UID)
		if err != nil {
			log.Printf("Error verifying user %s: %v", req.UID, err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		
		// Set custom claims
		claims := map[string]interface{}{
			"role": req.Role,
		}
		
		err = authClient.SetCustomUserClaims(context.Background(), req.UID, claims)
		if err != nil {
			log.Printf("Error setting custom claims for user %s: %v", req.UID, err)
			http.Error(w, "Failed to assign role", http.StatusInternalServerError)
			return
		}
		
		log.Printf("Role '%s' assigned to user %s (email: %s) by admin", req.Role, req.UID, user.Email)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Role assigned successfully",
			"uid":     req.UID,
			"role":    req.Role,
			"email":   user.Email,
		})
	}
}

// GetUserRoleHandler gets user's current role
func GetUserRoleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID")
		role := r.Context().Value("userRole")
		email := r.Context().Value("userEmail")
		
		// If context values are not set properly, try to get from UserContext
		if userID == nil || role == nil || email == nil {
			// Try to get from user context
			user, ok := r.Context().Value("user").(*UserContext)
			if !ok {
				http.Error(w, "User context not found", http.StatusUnauthorized)
				return
			}
			
			userIDStr := user.UID
			roleStr := user.Role
			emailStr := user.Email
			
			// Default to student if no role
			if roleStr == "" {
				roleStr = "student"
			}
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"uid":   userIDStr,
				"role":  roleStr,
				"email": emailStr,
			})
			return
		}
		
		userIDStr, _ := userID.(string)
		roleStr, _ := role.(string)
		emailStr, _ := email.(string)
		
		// Default to student if no role
		if roleStr == "" {
			roleStr = "student"
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"uid":   userIDStr,
			"role":  roleStr,
			"email": emailStr,
		})
	}
}

// Helper function to validate roles
func isValidRole(role string) bool {
	for _, validRole := range ValidRoles {
		if strings.ToLower(role) == validRole {
			return true
		}
	}
	return false
}

// UserContext struct for type checking (copied from middleware to avoid import cycle)
type UserContext struct {
	UID           string                 `json:"uid"`
	Email         string                 `json:"email"`
	EmailVerified bool                   `json:"email_verified"`
	Role          string                 `json:"role"`
	CollegeID     string                 `json:"college_id,omitempty"`
	Permissions   []string               `json:"permissions,omitempty"`
	CustomClaims  map[string]interface{} `json:"custom_claims,omitempty"`
}