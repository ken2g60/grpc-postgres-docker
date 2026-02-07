package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc_fintech/database"
	"github.com/grpc_fintech/internal/api/handlers"
	"github.com/grpc_fintech/internal/api/interceptors"
	"github.com/grpc_fintech/internal/models"
	mainapi "github.com/grpc_fintech/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func restapi_gatway() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	}))}

	err := mainapi.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50052", opts)
	if err != nil {
		log.Fatal("Failed to register gRPC-Gateway handler:", err)
	}

	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	log.Println("HTTPS Server is running on port: 8080...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start HTTP Server:", err)
	}
}

func serverGateway() {
	// connect to postgres
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors.ResponseTimeInterceptor))
	mainapi.RegisterUserServiceServer(s, &handlers.Server{})
	mainapi.RegisterWalletServiceServer(s, &handlers.Server{})

	fmt.Println("gRPC Server is running on port", ":50052")
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Printf("Error listen to port %v", err)
		return
	}

	reflection.Register(s)

	// listen to configuration
	err = s.Serve(lis)
	if err != nil {
		log.Printf("Error serve %s", err)
		return
	}

}

func main() {
	db := database.InitDb()
	// automigration listen to changes in the model
	migrations := database.Migrations{
		DB: db,
		Models: []interface{}{
			&models.User{},
			&models.Wallet{},
			&models.WalletHistory{},
			&models.Transaction{},
		},
	}
	database.RunMigrations(migrations)
	serverGateway()
	go restapi_gatway()
}
