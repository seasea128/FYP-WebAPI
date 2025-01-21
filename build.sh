#!/usr/bin/env bash

hash_str="`git rev-parse --short HEAD`"
version="`git describe --tags`"

if [ "$GOOS" = "windows" ]
then
		file_ext=".exe"
fi

go build -v -o ./bin/WebAPI-$version-$hash_str$file_ext #-ldflags "-X main.hash=$hash_str -X main.version=$version"
