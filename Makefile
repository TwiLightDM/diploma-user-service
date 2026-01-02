proto-gen:
	protoc --go_out=. --go-grpc_out=. proto/user_service.proto