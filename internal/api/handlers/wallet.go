package handlers

import (
	"context"
	"log"
	"time"

	"github.com/grpc_fintech/database"
	"github.com/grpc_fintech/internal/models"
	mainapi "github.com/grpc_fintech/proto/gen"
	"github.com/grpc_fintech/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Balance(ctx context.Context, req *mainapi.WalletIdRequest) (*mainapi.WalletResponse, error) {

	tokenString, err := utils.GetUserIdFromToken(ctx)
	userInfo, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token: "+err.Error())
	}

	wallet, err := models.GetWalletByUserID(ctx, database.Db, userInfo.UserId)
	if err != nil {
		return nil, status.Error(codes.Aborted, "user id not found")
	}

	return &mainapi.WalletResponse{Balance: float32(wallet.AvailableBalance)}, nil
}

func (s *Server) Deposit(ctx context.Context, req *mainapi.DepositRequest) (*mainapi.DepositResponse, error) {

	err := req.Validate()
	if err != nil {
		log.Printf("Validation failed: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	tokenString, err := utils.GetUserIdFromToken(ctx)
	userInfo, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token: "+err.Error())
	}

	wallet, err := models.GetWalletByUserID(ctx, database.Db, userInfo.UserId)
	if err != nil {
		return nil, status.Error(codes.Aborted, "user id not found")
	}

	wallet.UserID = userInfo.UserId
	wallet.AvailableBalance += req.Amount
	wallet.UpdatedAt = time.Now()
	if err := wallet.Save(database.Db); err != nil {
		return nil, status.Error(codes.Aborted, "Failed to update wallet balance")
	}

	return &mainapi.DepositResponse{Amount: wallet.AvailableBalance, Description: "withdrawal completed"}, nil
}

func (s *Server) Withdrawl(ctx context.Context, req *mainapi.WithdrawlRequest) (*mainapi.WithdrawlResponse, error) {

	err := req.Validate()
	if err != nil {
		log.Printf("Validation failed: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	tokenString, err := utils.GetUserIdFromToken(ctx)
	userInfo, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token: "+err.Error())
	}

	wallet, err := models.GetWalletByUserID(ctx, database.Db, userInfo.UserId)
	if err != nil {
		return nil, status.Error(codes.Aborted, "user id not found")
	}

	if float32(wallet.AvailableBalance) < req.Amount {
		return nil, status.Error(codes.Aborted, "wallet balance is less than requested amount")
	}

	wallet.UserID = userInfo.UserId
	wallet.AvailableBalance -= req.Amount
	wallet.UpdatedAt = time.Now()
	if err := wallet.Save(database.Db); err != nil {
		return nil, status.Error(codes.Aborted, "Failed to update wallet balance")
	}

	return &mainapi.WithdrawlResponse{Amount: wallet.AvailableBalance, Description: "withdrawal completed"}, nil
}

func (s *Server) TransactionHistory(ctx context.Context, req *mainapi.TransactionRequest) (*mainapi.TransactionHistoryResponse, error) {

	tokenString, err := utils.GetUserIdFromToken(ctx)
	userInfo, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token: "+err.Error())
	}

	transaction, err := models.TransactionHistory(ctx, database.Db, userInfo.UserId)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to fetch transaction history")
	}

	protoUsers := make([]*mainapi.TransactionRepeatedResponse, 0, len(*transaction))
	for _, model := range *transaction {
		protoUsers = append(protoUsers, &mainapi.TransactionRepeatedResponse{
			UserId:        model.UserID,
			PaymentMethod: model.PaymentMethod,
			Description:   model.Description,
			Amount:        float32(model.Amount),
		})
	}

	return &mainapi.TransactionHistoryResponse{Responses: protoUsers}, nil
}
