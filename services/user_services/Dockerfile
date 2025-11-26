FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /user_services .

FROM alpine:latest
WORKDIR /
COPY --from=builder /user_services /user_services
EXPOSE 8082
CMD ["/user_services"]