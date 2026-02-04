package interceptors

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ResponseTimeInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// Record the start time
	start := time.Now()

	// call the handler to proceed with the client request
	response, err := handler(ctx, req)

	// calculate the duration
	duration := time.Since(start)

	// log the request details with duration
	st, _ := status.FromError(err)
	fmt.Printf("Method : %s, Status: %d, Duration: %v\n", info.FullMethod, st.Code(), duration)

	md := metadata.Pairs("X-Response-Time", duration.String())
	grpc.SetHeader(ctx, md)

	return response, err
}
