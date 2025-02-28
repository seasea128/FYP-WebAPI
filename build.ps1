$hash_str = (git rev-parse --short HEAD).Trim()
$version = (git describe --tags).Trim()

$file_ext = if ($env:GOOS -eq "windows") { ".exe" } else { "" }

if ($args[0] -eq "docker") {
    $file_name = "WebAPI"
} else {
    $file_name = "WebAPI-$version-$hash_str$file_ext"
}

$PROTO_SRC_DIR = "./Protobuf-FYP/proto"
$DST_DIR = "."

protoc -I $PROTO_SRC_DIR --go_out $DST_DIR "$PROTO_SRC_DIR/data.proto"
protoc -I $PROTO_SRC_DIR --go_out $DST_DIR "$PROTO_SRC_DIR/timestamp.proto"

$bin_path = "./bin/$file_name"
go build -v -o $bin_path -ldflags "-X main.hash=$hash_str -X main.version=$version"
