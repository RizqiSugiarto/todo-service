#!/bin/bash

# Install dependencies
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate Stubs
protoc -I ./proto \
  --proto_path ./proto/ \
  --plugin=$(go env GOPATH)/bin/protoc-gen-go \
  --plugin=$(go env GOPATH)/bin/protoc-gen-go-grpc \
  --go_out ./stubs --go_opt paths=source_relative \
  --go-grpc_out ./stubs --go-grpc_opt paths=source_relative \
  ./proto/*/*.proto

