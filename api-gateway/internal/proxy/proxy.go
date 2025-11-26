package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func ForwardToUserService(w http.ResponseWriter, r *http.Request) {
    target, _ := url.Parse("http://user_services:8082") // User Service
    proxy := httputil.NewSingleHostReverseProxy(target)
    
    // User service uses /users endpoints (no /api prefix)
    // Strip /api/users and replace with just /users
    r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/users")
    r.URL.Path = "/users" + r.URL.Path
    if r.URL.Path == "/users" {
        r.URL.Path = "/users"
    }
    r.RequestURI = "" // Clear this so it's recalculated
    r.Header.Set("X-Forwarded-Host", r.Host)
    
    proxy.ServeHTTP(w, r)
}

func ForwardToBankingService(w http.ResponseWriter, r *http.Request) {
    target, _ := url.Parse("http://banking_services:8083") // Banking Service
    proxy := httputil.NewSingleHostReverseProxy(target)
    
    // When using chi Route, the full path is preserved
    // Remove the /api prefix and pass the request to the service
    r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
    if r.URL.Path == "" || r.URL.Path == "/banks" {
        r.URL.Path = "/api/banks"
    } else if strings.HasPrefix(r.URL.Path, "/banks/") {
        r.URL.Path = "/api" + r.URL.Path
    }
    r.RequestURI = "" // Clear this so it's recalculated
    r.Header.Set("X-Forwarded-Host", r.Host)
    
    proxy.ServeHTTP(w, r)
}

func ForwardToFinanceService(w http.ResponseWriter, r *http.Request) {
    target, _ := url.Parse("http://finance_services:8084") // Finance Service
    proxy := httputil.NewSingleHostReverseProxy(target)
    
    // Finance service expects /api/* paths
    // Strip /api/finance and replace with /api to match service routing
    r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/finance")
    r.URL.Path = "/api" + r.URL.Path
    if r.URL.Path == "/api" {
        r.URL.Path = "/api"
    }
    r.RequestURI = "" // Clear this so it's recalculated
    r.Header.Set("X-Forwarded-Host", r.Host)
    
    proxy.ServeHTTP(w, r)
}
