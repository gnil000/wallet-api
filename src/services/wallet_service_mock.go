package services

import (
	"context"
	"wallet-api/src/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type WalletServiceMock struct {
	mock.Mock
}

func (m *WalletServiceMock) GetWalletByID(ctx context.Context, id uuid.UUID) (models.GetBalanceResponse, *models.ErrorResponse) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.GetBalanceResponse), args.Get(1).(*models.ErrorResponse)
}

func (m *WalletServiceMock) ChangeWalletBalance(ctx context.Context, req models.ChangeBalanceRequest) *models.ErrorResponse {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*models.ErrorResponse)
}

func (m *WalletServiceMock) GetWallets(ctx context.Context) ([]models.GetWalletsResponse, *models.ErrorResponse) {
	args := m.Called(ctx)
	return args.Get(0).([]models.GetWalletsResponse), args.Get(1).(*models.ErrorResponse)
}
