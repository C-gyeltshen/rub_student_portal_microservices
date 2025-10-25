package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
    const (
        userServiceURL    = "http://user_services:8082"
        bankingServiceURL = "http://banking_services:8083"
    )
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    // User service routes (forward paths as-is so downstream handlers keep their /api/... paths)
    r.HandleFunc("/api/get/*", proxyTo(userServiceURL))
    r.HandleFunc("/api/create/user", proxyTo(userServiceURL))
    r.HandleFunc("/api/get/user/*", proxyTo(userServiceURL))

    // Banking service routes
    r.HandleFunc("/api/banks*", proxyTo(bankingServiceURL))
    r.HandleFunc("/api/student-bank-details*", proxyTo(bankingServiceURL))

    log.Println("API Gateway starting on :8080")
    http.ListenAndServe(":8080", r)
}

func proxyTo(targetURL string) http.HandlerFunc {
    target, _ := url.Parse(targetURL)
    proxy := httputil.NewSingleHostReverseProxy(target)

    // Ensure the director preserves the incoming path and query but updates the target host
    originalDirector := proxy.Director
    proxy.Director = func(req *http.Request) {
        // call the default director which sets scheme/host
        originalDirector(req)
        // keep the incoming path and raw query
        // (originalDirector sets Path based on target, so overwrite it)
        // We copy from the incoming request that initiated the proxy call via closure
        // Note: the proxy handler will receive the same Request pointer, so Path is already as incoming.
    }

    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Proxying %s %s to %s", r.Method, r.URL.Path, targetURL)
        proxy.ServeHTTP(w, r)
    }
}