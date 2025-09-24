package repositories

import (
	"context"
	"wallet-api/src/database/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type WalletRepoMock struct {
	mock.Mock
}

func NewWalletRepoMock() *WalletRepoMock {
	return &WalletRepoMock{}
}

func (m *WalletRepoMock) FindByID(ctx context.Context, id uuid.UUID) (entities.Wallet, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.Wallet), args.Error(1)
}

func (m *WalletRepoMock) DepositUpdate(ctx context.Context, id uuid.UUID, amount int64) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *WalletRepoMock) WithdrawUpdate(ctx context.Context, id uuid.UUID, amount int64) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *WalletRepoMock) GetWallets(ctx context.Context) ([]entities.Wallet, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Wallet), args.Error(0)
}
