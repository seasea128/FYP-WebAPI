# FinalYearProject API Server

## Requirements
 - Go >= 1.23.2

## Building

### Linux
``` bash
go get -u
./build.sh
```

### Windows
``` powershell
go get -u
./build.ps1
```

## Configuration
A configuration file will be created on initial launch of the server. The configuration file is describe below.

| Field Name       | Type   | Default value                              |
|------------------|--------|--------------------------------------------|
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
