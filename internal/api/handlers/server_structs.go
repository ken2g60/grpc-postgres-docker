package handlers

import mainapi "github.com/grpc_fintech/proto/gen"

type Server struct {
	mainapi.UnimplementedUserServiceServer
	mainapi.UnimplementedWalletServiceServer
}
