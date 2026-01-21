proto-gen:
	protoc --go_out=. --go-grpc_out=. proto/user_service.proto
	protoc --go_out=. --go-grpc_out=. proto/group_service.proto
	protoc --go_out=. --go-grpc_out=. proto/group_member_service.proto