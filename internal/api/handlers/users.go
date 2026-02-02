package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/grpc_fintech/database"
	"github.com/grpc_fintech/internal/models"
	mainapi "github.com/grpc_fintech/proto/gen"
	"github.com/grpc_fintech/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *mainapi.CreateUserRequest) (*mainapi.UserResponse, error) {

	// validate email
	user, err := models.FindUser(ctx, database.Db, req.GetEmail())
	if err != nil {
		log.Println("error fetching email data")
		return nil, status.Error(codes.Aborted, "error fetching data")
	}

	if user.Email != "" {
		return nil, status.Error(codes.AlreadyExists, "account already exists")
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

	return &mainapi.UserResponse{Email: req.Email, FirstName: req.FirstName, LastName: req.LastName}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *mainapi.LoginRequest) (*mainapi.LoginResponse, error) {

	user, err := models.FindUser(ctx, database.Db, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user.Email == "" {
		return nil, status.Error(codes.Aborted, "email does not exits")
	}

	err = utils.VerifyPassword(req.GetPassword(), user.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "incorrect username/password")
	}

	tokenString, err := utils.SignToken(user.UUID, user.Email)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "could not create token.")
	}

	return &mainapi.LoginResponse{Token: tokenString, Status: true}, nil
}

func (s *Server) UserProfile(ctx context.Context, req *mainapi.UserIdRequest) (*mainapi.UserResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata found")
	}

	val, ok := md["authorization"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized Access")
	}

	tokenString := strings.TrimPrefix(val[0], "Bearer ")
	if tokenString == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized Access")
	}

	userInfo, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token: "+err.Error())
	}

	user, err := models.FindUser(ctx, database.Db, userInfo.UserId)
	if err != nil {
		log.Println("error fetching email data")
		return nil, status.Error(codes.Aborted, "error fetching data")
	}

	return &mainapi.UserResponse{FirstName: user.First_name, LastName: user.Last_name, Email: user.Email}, nil
}

func (s *Server) UpdateProfile(ctx context.Context, req *mainapi.UserID) (*mainapi.UserResponse, error) {
	fmt.Println("implement the update function")
	return &mainapi.UserResponse{}, nil
}

func (s *Server) DeactivateAccount(ctx context.Context, req *mainapi.UserID) (*mainapi.DeactivateResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata found")
	}

	val, ok := md["authorization"]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized Access")
	}

	tokenString := strings.TrimPrefix(val[0], "Bearer ")
	if tokenString == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized Access")
	}

	userInfo, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token: "+err.Error())
	}

	user, err := models.DeactivateAccount(ctx, database.Db, userInfo.UserId)
	if err != nil {
		return nil, status.Error(codes.Aborted, "error deactivating account")
	}

	return &mainapi.DeactivateResponse{Status: user.Status}, nil
}
