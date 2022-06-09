protoc --proto_path=../api --go_out=../pkg/proto --go-grpc_out=../pkg/proto --go_opt=paths=source_relative service.proto
