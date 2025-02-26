FROM golang:1.24.0-bookworm

WORKDIR /usr/src/api

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./ ./

RUN apt update
RUN apt -y upgrade
RUN apt -y install protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

RUN ./build.sh docker

RUN apt -y remove protobuf-compiler

RUN useradd api
RUN chown api:api /usr/src/api
USER api

CMD ["./bin/WebAPI"]
