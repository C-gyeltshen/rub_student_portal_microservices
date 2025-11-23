package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ForwardToUserService(w http.ResponseWriter, r *http.Request) {
    target, _ := url.Parse("http://localhost:8082") // User Service
    proxy := httputil.NewSingleHostReverseProxy(target)
    proxy.ServeHTTP(w, r)
}