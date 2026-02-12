package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) ApplyOperation(
	ctx context.Context,
	walletID string,
	amount int64,
) (int64, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var balance int64

	err = tx.QueryRowContext(
		ctx,
		`SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`,
		walletID,
	).Scan(&balance)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("wallet not found")
	}
	if err != nil {
		if isInvalidUUIDError(err) {
			return 0, fmt.Errorf("invalid wallet id format")
		}
		return 0, err
	}

	newBalance := balance + amount
	if newBalance < 0 {
		return 0, fmt.Errorf("insufficient funds")
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE wallets SET balance = $1 WHERE id = $2`,
		newBalance,
		walletID,
	)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newBalance, nil
}

func (r *WalletRepository) GetBalance(
	ctx context.Context,
	walletID string,
) (int64, error) {

	var balance int64

	err := r.db.QueryRowContext(
		ctx,
		`SELECT balance FROM wallets WHERE id = $1`,
		walletID,
	).Scan(&balance)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("wallet not found")
	}
	if err != nil {
		if isInvalidUUIDError(err) {
			return 0, fmt.Errorf("invalid wallet id format")
		}
		return 0, err
	}

	return balance, nil
}

func isInvalidUUIDError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return contains(errMsg, "invalid input syntax for type uuid") ||
		contains(errMsg, "invalid UUID")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
