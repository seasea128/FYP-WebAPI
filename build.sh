#!/usr/bin/env bash

hash_str="`git rev-parse --short HEAD`"
version="`git describe --tags`"

if [ "$GOOS" = "windows" ]
then
	file_ext=".exe"
fi

if [ "$1" = "docker" ]
then
	file_name="WebAPI"
	# export CGO_ENABLED=0
else
	file_name="WebAPI-$version-$hash_str$file_ext"
fi

PROTO_SRC_DIR="./Protobuf-FYP/proto"
DST_DIR="."

protoc -I=$PROTO_SRC_DIR --go_out=$DST_DIR $PROTO_SRC_DIR/data.proto
protoc -I=$PROTO_SRC_DIR --go_out=$DST_DIR $PROTO_SRC_DIR/timestamp.proto

go build -v -o ./bin/$file_name -ldflags "-X main.hash=$hash_str -X main.version=$version"
