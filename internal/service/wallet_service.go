package service

import (
	"context"
	"fmt"
)

type WalletRepository interface {
	ApplyOperation(ctx context.Context, walletID string, amount int64) (int64, error)
	GetBalance(ctx context.Context, walletID string) (int64, error)
}

type WalletService struct {
	repo WalletRepository
}

func NewWalletService(repo WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) Apply(
	ctx context.Context,
	walletID string,
	operation string,
	amount int64,
) (int64, error) {

	if amount <= 0 {
		return 0, fmt.Errorf("amount must be positive")
	}

	switch operation {
	case "DEPOSIT":
		return s.repo.ApplyOperation(ctx, walletID, amount)
	case "WITHDRAW":
		return s.repo.ApplyOperation(ctx, walletID, -amount)
	default:
		return 0, fmt.Errorf("unknown operation type")
	}
}

func (s *WalletService) GetBalance(
	ctx context.Context,
	walletID string,
) (int64, error) {
	return s.repo.GetBalance(ctx, walletID)
}
