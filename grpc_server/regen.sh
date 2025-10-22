#!/bin/sh

build_tools() {
    mkdir -p bin
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
}

generate() {
    mkdir -p proto/generated_files
    protoc -I=proto --go_out=proto/generated_files --go-grpc_out=proto/generated_files \
        --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/*.proto

}

export GOBIN=$PWD/bin
export PATH=$GOBIN:$PATH
build_tools
generate