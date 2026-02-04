package main

import (
	"fmt"
	"log"
	"net"

	"github.com/grpc_fintech/database"
	"github.com/grpc_fintech/internal/api/handlers"
	"github.com/grpc_fintech/internal/api/interceptors"
	"github.com/grpc_fintech/internal/models"
	mainapi "github.com/grpc_fintech/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// connect to postgres
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
