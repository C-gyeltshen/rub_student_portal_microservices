package proxy

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// UserContextKey is the type used for context keys
type UserContextKey string

const (
	// UserContextKeyValue is the context key for user information
	UserContextKeyValue UserContextKey = "user"
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

// GetUserFromContext extracts user information from request context
func GetUserFromContext(ctx context.Context) (*UserContext, bool) {
	user, ok := ctx.Value(UserContextKeyValue).(*UserContext)
	return user, ok
}

func ForwardToUserService(w http.ResponseWriter, r *http.Request) {
    target, _ := url.Parse("http://user_services:8082") // User Service
    proxy := httputil.NewSingleHostReverseProxy(target)
    
    // Strip /api prefix from the path
    r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
    
    // Forward user context headers if user is authenticated
    if user, ok := GetUserFromContext(r.Context()); ok {
        r.Header.Set("X-User-ID", user.UID)
        r.Header.Set("X-User-Email", user.Email)
        r.Header.Set("X-User-Role", user.Role)
        if user.CollegeID != "" {
            r.Header.Set("X-User-College-ID", user.CollegeID)
        }
        
        // Add permissions as comma-separated header
        if len(user.Permissions) > 0 {
            r.Header.Set("X-User-Permissions", strings.Join(user.Permissions, ","))
        }
        
        log.Printf("Forwarding request to user service with user context: UID=%s, Role=%s", user.UID, user.Role)
    }
    
    r.Header.Set("X-Forwarded-Host", r.Host)
    
    // Remove the Authorization header to prevent downstream services from processing it
    r.Header.Del("Authorization")
    
    proxy.ServeHTTP(w, r)
}

func ForwardToBankingService(w http.ResponseWriter, r *http.Request) {
    target, _ := url.Parse("http://banking_services:8083") // Banking Service
    proxy := httputil.NewSingleHostReverseProxy(target)
    
    // Strip /api prefix from the path
    r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
    
    // Forward user context headers if user is authenticated
    if user, ok := GetUserFromContext(r.Context()); ok {
        r.Header.Set("X-User-ID", user.UID)
        r.Header.Set("X-User-Email", user.Email)
        r.Header.Set("X-User-Role", user.Role)
        if user.CollegeID != "" {
            r.Header.Set("X-User-College-ID", user.CollegeID)
        }
        
        // Add permissions as comma-separated header
        if len(user.Permissions) > 0 {
            r.Header.Set("X-User-Permissions", strings.Join(user.Permissions, ","))
        }
        
        log.Printf("Forwarding request to banking service with user context: UID=%s, Role=%s", user.UID, user.Role)
    }
    
    r.Header.Set("X-Forwarded-Host", r.Host)
    
    // Remove the Authorization header to prevent downstream services from processing it
    r.Header.Del("Authorization")
    
    proxy.ServeHTTP(w, r)
}