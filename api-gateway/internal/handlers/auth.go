package handlers


import (
	"encoding/json"
	"net/http"
	"log"

	"api-gateway/internal/middleware"
)

// DashboardResponse represents the dashboard data
type DashboardResponse struct {
	Message string      `json:"message"`
	User    interface{} `json:"user,omitempty"`
}

// HealthResponse represents health check data
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// Dashboard handler - public endpoint for login/dashboard access
func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Check if user is authenticated (optional)
	if user, ok := middleware.GetUserFromContext(r.Context()); ok {
		response := DashboardResponse{
			Message: "Welcome to RUB Student Portal Dashboard",
			User: map[string]interface{}{
				"uid":   user.UID,
				"email": user.Email,
				"role":  user.Role,
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Not authenticated - show public dashboard
	response := DashboardResponse{
		Message: "Welcome to RUB Student Portal - Please login to continue",
	}
	json.NewEncoder(w).Encode(response)
}

// Health check endpoint - public
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: "2025-11-26T00:00:00Z",
		Version:   "1.0.0",
	}
	
	json.NewEncoder(w).Encode(response)
}

// UserProfile handler - requires authentication
func UserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "AUTH_001")
		return
	}

	profile := map[string]interface{}{
		"uid":            user.UID,
		"email":          user.Email,
		"email_verified": user.EmailVerified,
		"role":           user.Role,
		"college_id":     user.CollegeID,
		"permissions":    user.Permissions,
	}

	response := map[string]interface{}{
		"success": true,
		"data":    profile,
		"message": "User profile retrieved successfully",
	}

	json.NewEncoder(w).Encode(response)
}

// AdminDashboard handler - requires admin role
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "AUTH_001")
		return
	}

	adminData := map[string]interface{}{
		"total_users":    1250,
		"total_students": 1100,
		"total_colleges": 8,
		"active_budgets": 15,
		"pending_requests": 23,
		"user_info": map[string]interface{}{
			"uid":   user.UID,
			"email": user.Email,
			"role":  user.Role,
		},
	}

	response := map[string]interface{}{
		"success": true,
		"data":    adminData,
		"message": "Admin dashboard data retrieved successfully",
	}

	json.NewEncoder(w).Encode(response)
}

// FinanceOfficerDashboard handler - requires finance officer or admin role
func FinanceOfficerDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "AUTH_001")
		return
	}

	financeData := map[string]interface{}{
		"total_budget":      500000.00,
		"allocated_budget":  350000.00,
		"remaining_budget":  150000.00,
		"pending_expenses":  25,
		"approved_expenses": 145,
		"college_id":        user.CollegeID,
		"user_info": map[string]interface{}{
			"uid":   user.UID,
			"email": user.Email,
			"role":  user.Role,
		},
	}

	response := map[string]interface{}{
		"success": true,
		"data":    financeData,
		"message": "Finance officer dashboard data retrieved successfully",
	}

	json.NewEncoder(w).Encode(response)
}

// StudentDashboard handler - requires student role or higher
func StudentDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "AUTH_001")
		return
	}

	studentData := map[string]interface{}{
		"student_id":        "STU2025001",
		"college_id":        user.CollegeID,
		"enrollment_status": "Active",
		"current_semester":  "Fall 2025",
		"stipend_status":    "Pending",
		"bank_details":      "Configured",
		"user_info": map[string]interface{}{
			"uid":   user.UID,
			"email": user.Email,
			"role":  user.Role,
		},
	}

	response := map[string]interface{}{
		"success": true,
		"data":    studentData,
		"message": "Student dashboard data retrieved successfully",
	}

	json.NewEncoder(w).Encode(response)
}

// LoginInstructions provides instructions for Firebase authentication
func LoginInstructions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	instructions := map[string]interface{}{
		"message": "Firebase Authentication Required",
		"instructions": []string{
			"1. Authenticate with Firebase using your email and password",
			"2. Obtain the Firebase ID token from your client application",
			"3. Include the token in the Authorization header as 'Bearer <token>'",
			"4. Access protected endpoints with proper roles and permissions",
		},
		"endpoints": map[string]interface{}{
			"dashboard":      "/dashboard (public)",
			"health":         "/health (public)",
			"user_profile":   "/profile (authenticated users)",
			"admin":          "/admin/dashboard (admin only)",
			"finance":        "/finance/dashboard (finance officer + admin)",
			"student":        "/student/dashboard (all authenticated users)",
		},
		"roles": []string{
			"admin",
			"finance_officer", 
			"student",
		},
	}

	response := map[string]interface{}{
		"success": true,
		"data":    instructions,
	}

	log.Printf("Login instructions requested from %s", r.RemoteAddr)
	json.NewEncoder(w).Encode(response)
}