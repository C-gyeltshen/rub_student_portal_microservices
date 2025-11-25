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
    
    // Strip /api/users prefix from the path, keeping remaining path
    r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/users")
    if r.URL.Path == "" {
        r.URL.Path = "/"
    }
    r.RequestURI = "" // Clear this so it's recalculated
    r.Header.Set("X-Forwarded-Host", r.Host)
    
    proxy.ServeHTTP(w, r)
}

func ForwardToBankingService(w http.ResponseWriter, r *http.Request) {
    target, _ := url.Parse("http://banking_services:8083") // Banking Service
    proxy := httputil.NewSingleHostReverseProxy(target)
    
    // Strip /api/banks prefix from the path, keeping remaining path
    r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/banks")
    if r.URL.Path == "" {
        r.URL.Path = "/"
    }
    r.RequestURI = "" // Clear this so it's recalculated
    r.Header.Set("X-Forwarded-Host", r.Host)
    
    proxy.ServeHTTP(w, r)
}

func ForwardToFinanceService(w http.ResponseWriter, r *http.Request) {
    target, _ := url.Parse("http://finance_services:8084") // Finance Service
    proxy := httputil.NewSingleHostReverseProxy(target)
    
    // When using Chi's Mount, the path is still the original (/api/finance/...)
    // Strip /api/finance and then add /api back (finance service expects /api/...)
    r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/finance")
    r.URL.Path = "/api" + r.URL.Path
    if r.URL.Path == "/api" {
        r.URL.Path = "/api/"
    }
    r.RequestURI = "" // Clear this so it's recalculated
    r.Header.Set("X-Forwarded-Host", r.Host)
    
    proxy.ServeHTTP(w, r)
}
