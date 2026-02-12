package test

import (
	"context"
	"testing"

	"wallet-api/internal/service"
)

type mockRepo struct{}

func (m *mockRepo) ApplyOperation(ctx context.Context, id string, amount int64) (int64, error) {
	return 100 + amount, nil
}

func (m *mockRepo) GetBalance(ctx context.Context, id string) (int64, error) {
	return 100, nil
}

func TestDeposit(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewWalletService(repo)

	balance, err := svc.Apply(context.Background(), "id", "DEPOSIT", 50)
	if err != nil {
		t.Fatal(err)
	}

	if balance != 150 {
		t.Fatalf("expected 150, got %d", balance)
	}
}

func TestGetBalance(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewWalletService(repo)

	balance, err := svc.GetBalance(context.Background(), "id")
	if err != nil {
		t.Fatal(err)
	}

	if balance != 100 {
		t.Fatalf("expected 100, got %d", balance)
	}
}
