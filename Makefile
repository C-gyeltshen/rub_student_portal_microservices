.PHONY: proto-user proto-banking proto-all

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