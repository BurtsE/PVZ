package generated

//go:generate oapi-codegen -generate types -package openapi -o ./openapi/types.go ../api/swagger.yaml
//go:generate oapi-codegen -generate gin -package openapi -o ./openapi/server.go ../api/swagger.yaml
//go:generate oapi-codegen -generate client -package openapi -o ./openapi/client.go ../api/swagger.yaml
//go:generate oapi-codegen -generate spec -package openapi -o ./openapi/spec.go ../api/swagger.yaml
