```bash
api-gateway/
│
├── cmd/
│   └── main.go               # Entry point (bootstraps server)
│
├── internal/
│   ├── config/
│   │   └── config.go         # Environment variables, service URLs, etc.
│   │
│   ├── middleware/
│   │   ├── auth.go           # JWT / API key verification
│   │   ├── logger.go         # Request logging
│   │   └── rate_limit.go     # Optional: rate limiting
│   │
│   ├── router/
│   │   └── router.go         # Route definitions using chi/gin/fiber
│   │
│   ├── proxy/
│   │   ├── proxy.go          # Reverse proxy logic (e.g., httputil.ReverseProxy)
│   │   └── grpc_client.go    # For gRPC service calls
│   │
│   ├── validators/
│   │   └── request_validator.go # Parameter validation / schema validation
│   │
│   ├── services/
│   │   ├── user_service.go   # HTTP handler or gRPC client wrapper
│   │   ├── order_service.go
│   │   └── payment_service.go
│   │
│   ├── utils/
│   │   └── response.go       # Common response/error format
│   │
│   └── handlers/
│       └── gateway_handler.go # Request entrypoint handler
│
├── pkg/
│   └── logger/
│       └── logger.go         # Centralized logging utilities
│
├── go.mod
└── go.sum
```