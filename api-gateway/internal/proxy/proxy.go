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
    
    // Strip /api prefix from the path
    r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
    r.Header.Set("X-Forwarded-Host", r.Host)
    
    proxy.ServeHTTP(w, r)
}

func ForwardToBankingService(w http.ResponseWriter, r *http.Request) {
    target, _ := url.Parse("http://banking_services:8083") // Banking Service
    proxy := httputil.NewSingleHostReverseProxy(target)
    
    // Strip /api prefix from the path
    r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
    r.Header.Set("X-Forwarded-Host", r.Host)
    
    proxy.ServeHTTP(w, r)
}