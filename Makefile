.PHONY: proto-user proto-banking proto-all test test-coverage test-banking test-user test-all test-coverage-check

proto-user:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user/user.proto

proto-banking:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/banking/banking.proto

proto-all: proto-user proto-banking

clean-proto:
	find proto -name "*.pb.go" -delete

# Test targets
test:
	@echo "Running all tests..."
	@cd services/banking_services && go test ./... -v
	@cd services/user_services && go test ./... -v

test-coverage:
	@echo "Running tests with coverage..."
	@cd services/banking_services && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
	@cd services/user_services && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage reports generated: services/*/coverage.html"

test-banking:
	@echo "Testing banking service..."
	@cd services/banking_services && go test ./... -v -cover

test-user:
	@echo "Testing user service..."
	@cd services/user_services && go test ./... -v -cover

test-all: test-coverage
	@echo "All tests completed with coverage"

# Coverage threshold check
test-coverage-check:
	@echo "Checking coverage thresholds..."
	@cd services/banking_services && go test ./... -coverprofile=coverage.out
	@cd services/banking_services && go tool cover -func=coverage.out | grep total | awk '{if ($$3+0 < 80.0) {print "❌ Banking service coverage below 80%: " $$3; exit 1} else {print "✅ Banking service coverage: " $$3}}'
	@cd services/user_services && go test ./... -coverprofile=coverage.out
	@cd services/user_services && go tool cover -func=coverage.out | grep total | awk '{if ($$3+0 < 80.0) {print "❌ User service coverage below 80%: " $$3; exit 1} else {print "✅ User service coverage: " $$3}}'
	@echo "✓ All services meet 80% coverage threshold"