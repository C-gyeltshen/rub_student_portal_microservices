FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /student_management_service .

FROM alpine:latest
WORKDIR /
COPY --from=builder /student_management_service /student_management_service
EXPOSE 8084
CMD ["/student_management_service"]
