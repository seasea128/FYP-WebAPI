# FinalYearProject API Server

## Requirements

- Go >= 1.23.2
- protoc >= 29.2 (Can be installed via winget or choco on Windows)

### SQLite3 support

- GCC >= 14 (Tested on GCC 14.2.1, older version should work as well)

## Building

### Linux

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
git submodule update --init --recursive .
go get -u
./build.sh
```

### Windows (Powershell)

```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
git submodule update --init --recursive .
go get -u
./build.ps1
```

## Deployment

### Docker

#### Remote deployment

Clone the repo to server or connect to server's Docker instance with docker context

```bash
docker context create <context-name> --docker "host=ssh://<username>@<hostname>"
docker context use <context-name>
```

Then, the project can be deployed with docker compose.

```bash
docker compose up -d
```

This might take a while to build the image, as go-sqlite are used and it's a package that is known to take quite a long time to build because of CGO.

There is 3 containers being deployed with this compose file, first is a postgres database that will be used to store all of the data, second is a grafana instance that will be used as the dashboard, and finally the API itself.

## Configuration

A configuration file will be created on initial launch of the server. The configuration file is describe below.

| Field Name       | Type   | Default value                              |
| ---------------- | ------ | ------------------------------------------ |
| DBType           | DBType | Default: sqlite, can be postgres or sqlite |
| ConnectionString | String | test.db                                    |
| HTTPPort         | int    | 8080                                       |
| MQTTPort         | int    | 1433                                       |

## HTTP API Documentation

The documentation for the API is generated inside the server itself. It can be accessed by visiting /docs.

```
http://localhost:8080/docs
```

`
