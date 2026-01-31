package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/grpc_fintech/database"
	"github.com/grpc_fintech/internal/api/handlers"
	"github.com/grpc_fintech/internal/models"
	mainapi "github.com/grpc_fintech/proto/gen"

	"google.golang.org/grpc"
)

func main() {

	// connect to postgres
	db := database.InitDb()
	// automigration listen to changes in the model
	migrations := database.Migrations{
		DB: db,
		Models: []interface{}{
			&models.User{},
		},
	}
	database.RunMigrations(migrations)

	s := grpc.NewServer()
	mainapi.RegisterUserServiceServer(s, &handlers.Server{})

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":50021"
	}
	fmt.Println("gRPC Server is running on port", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("Error listen to port %v", err)
		return
	}

	// listen to configuration
	err = s.Serve(lis)
	if err != nil {
		log.Printf("Error serve %s", err)
		return
	}

}
