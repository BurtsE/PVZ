package generated

//go:generate oapi-codegen -generate types -package openapi -o ./openapi/types.go ../api/swagger.yaml
//go:generate oapi-codegen -generate gin -package openapi -o ./openapi/server.go ../api/swagger.yaml
//go:generate oapi-codegen -generate client -package openapi -o ./openapi/client.go ../api/swagger.yaml
//go:generate oapi-codegen -generate spec -package openapi -o ./openapi/spec.go ../api/swagger.yaml
//go:generate protoc --proto_path=../api --go_out=./protobuf/ --go_opt=paths=source_relative --go-grpc_out=./protobuf/ --go-grpc_opt=paths=source_relative pvz.proto
