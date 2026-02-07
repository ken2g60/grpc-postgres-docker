go get google.golang.org/grpc
go get github.com/grpc-ecosystem/grpc-gateway/v2/runtime@v2.21.0

go install \
 github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
 github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 (exclude)

go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

wget -P proto/google/api https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto (assuming you are at project folder root)

curl -o annotations.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto
curl -o http.proto https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto

protoc -I=proto --go_out=. --go-grpc_out=. --grpc-gateway_out=. --validate_out=lang=go:. proto/users.proto proto/wallet.proto
