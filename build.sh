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
else
	file_name="WebAPI-$version-$hash_str$file_ext"
	export CGO_ENABLED=0
fi

PROTO_SRC_DIR="./Protobuf-FYP/proto"
DST_DIR="."

protoc -I=$PROTO_SRC_DIR --go_out=$DST_DIR $PROTO_SRC_DIR/*

go build -o ./bin/$file_name -ldflags "-X main.hash=$hash_str -X main.version=$version"
