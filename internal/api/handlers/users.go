package handlers

import (
	"context"
	"time"

	"github.com/grpc_fintech/database"
	"github.com/grpc_fintech/internal/models"
	mainapi "github.com/grpc_fintech/proto/gen"
	"github.com/grpc_fintech/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *mainapi.CreateUserRequest) (*mainapi.UserResponse, error) {

	// validate email
	user, err := models.FindUser(ctx, database.Db, req.Email)
	if err != nil {
		return &mainapi.UserResponse{}, nil
	}

	if user.Email != "" {
		return &mainapi.UserResponse{}, nil
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return &mainapi.UserResponse{}, status.Error(codes.Aborted, err.Error())
	}

	user_data := models.User{
		First_name: req.FirstName,
		Last_name:  req.LastName,
		Email:      req.Email,
		Password:   hashedPassword,
		CreatedAt:  time.Now(),
	}

	err = models.CreateUser(ctx, database.Db, &user_data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &mainapi.UserResponse{Email: req.Email}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *mainapi.LoginRequest) (*mainapi.LoginResponse, error) {

	user, err := models.FindUser(ctx, database.Db, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user.Email != "" {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = utils.VerifyPassword(req.GetPassword(), user.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "incorrect username/password")
	}

	tokenString, err := utils.SignToken(user.UUID, user.First_name, user.Last_name)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "could not create token.")
	}

	return &mainapi.LoginResponse{Token: tokenString, Status: true}, nil
}
