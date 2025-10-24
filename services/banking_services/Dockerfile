FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /banking_services .

FROM alpine:latest
WORKDIR /
COPY --from=builder /banking_services /banking_services
EXPOSE 8083
CMD ["/banking_services"]
