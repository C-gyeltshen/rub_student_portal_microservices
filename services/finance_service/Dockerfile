# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install protoc and grpc tools
RUN apk add --no-cache protobuf protobuf-dev git build-base

# Install Go protobuf and gRPC plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .

# Generate gRPC code from proto files
RUN mkdir -p pkg/pb && \
    protoc --go_out=pkg/pb --go-grpc_out=pkg/pb -I proto proto/stipend.proto proto/deduction.proto

# Download Go dependencies
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o finance_service .

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/finance_service .

# Expose both REST API and gRPC ports
EXPOSE 8084 50051

# Environment variables
ENV PORT=8084 \
    GRPC_PORT=50051

CMD ["./finance_service"]
