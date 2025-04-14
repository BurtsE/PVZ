FROM golang:1.23-alpine AS builder

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY cmd cmd
COPY api api
COPY internal internal
COPY pkg pkg
COPY generated/generate.go generated/generate.go
COPY .env .env

RUN mkdir generated/openapi
RUN mkdir generated/protobuf

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
RUN GOOS=linux go generate ./generated
RUN GOOS=linux go build -o app ./cmd

FROM golang:1.23-alpine
WORKDIR /root/


COPY --from=builder /build/app .

EXPOSE 8080

CMD ["./app"]