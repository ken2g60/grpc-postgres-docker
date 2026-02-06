protoc -I=proto --go_out=. --go-grpc_out=. --validate_out=lang=go:. proto/users.proto proto/wallet.proto
