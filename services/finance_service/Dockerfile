# Build stage
FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o finance_service .

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/finance_service .

EXPOSE 8084

CMD ["./finance_service"]
