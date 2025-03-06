package generated

//go:generate protoc --proto_path=../api --go_out=./protobuf/ --go_opt=paths=source_relative --go-grpc_out=./protobuf/ --go-grpc_opt=paths=source_relative api/proto/rates.proto api/proto/rates.proto
