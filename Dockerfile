# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY . ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the gRPC server
RUN go build -o bin/server ./cmd/grpcapi/binaries/serve.go

# Build the gRPC-Gateway
RUN go build -o bin/gateway ./cmd/grpcapi/binaries/serve.go

# Run Stage for gRPC Server
FROM alpine:latest AS server

WORKDIR /app

COPY --from=builder /app/bin/server .

EXPOSE 50051

ENTRYPOINT ["./server"]

# Run Stage for gRPC-Gateway
FROM alpine:latest AS gateway

WORKDIR /app

COPY --from=builder /app/bin/gateway .

EXPOSE 8080

ENTRYPOINT ["./gateway"]
