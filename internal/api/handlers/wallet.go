package handlers

import (
	"context"

	"github.com/grpc_fintech/database"
	"github.com/grpc_fintech/internal/models"
	mainapi "github.com/grpc_fintech/proto/gen"
)

func (s *Server) Balance(ctx context.Context, req *mainapi.WalletIdRequest) (*mainapi.WalletResponse, error) {
	wallet, err := models.GetWalletByUserID(ctx, database.Db, "user-id")
	if err != nil {
	}

	return &mainapi.WalletResponse{Balance: float32(wallet.AvailableBalance)}, nil
}

func (s *Server) Deposit(ctx context.Context, req *mainapi.DepositRequest) (*mainapi.DepositResponse, error) {
	return &mainapi.DepositResponse{Amount: 0, Description: ""}, nil
}

func (s *Server) Withdrawl(ctx context.Context, req *mainapi.WithdrawlRequest) (*mainapi.WithdrawlResponse, error) {
	return &mainapi.WithdrawlResponse{}, nil
}
