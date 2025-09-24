package services_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"wallet-api/src/database/entities"
	"wallet-api/src/database/repositories"
	"wallet-api/src/models"
	"wallet-api/src/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWalletService_GetWalletByID_Success(t *testing.T) {
	mockRepo := new(repositories.WalletRepoMock)
	svc := services.NewWalletService(mockRepo)

	ctx := context.Background()
	validID := uuid.New()
	t.Run("success", func(t *testing.T) {
		entity := entities.Wallet{
			ID:      validID,
			Balance: 1000,
		}
		mockRepo.On("FindByID", ctx, validID).Return(entity, nil)

		wallet, errResp := svc.GetWalletByID(ctx, validID)
		assert.Nil(t, errResp)
		assert.Equal(t, entity.Balance, wallet.Balance)

		mockRepo.AssertExpectations(t)
	})
}

func TestWalletService_GetWalletByID(t *testing.T) {
	mockRepo := new(repositories.WalletRepoMock)
	svc := services.NewWalletService(mockRepo)

	ctx := context.Background()
	validID := uuid.New()

	t.Run("wallet not found", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, validID).Return(entities.Wallet{}, repositories.ErrWalletNotFound)

		wallet, errResp := svc.GetWalletByID(ctx, validID)
		assert.Empty(t, wallet)
		assert.Equal(t, http.StatusNotFound, errResp.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("internal error on find", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, validID).Return(entities.Wallet{}, errors.New("db error"))

		wallet, errResp := svc.GetWalletByID(ctx, validID)
		assert.Empty(t, wallet)
		assert.Equal(t, http.StatusNotFound, errResp.Code)

		mockRepo.AssertExpectations(t)
	})
}

func TestWalletService_ChangeWalletBalance(t *testing.T) {
	mockRepo := new(repositories.WalletRepoMock)
	svc := services.NewWalletService(mockRepo)
	ctx := context.Background()
	walletID := uuid.New()

	t.Run("deposit success", func(t *testing.T) {
		req := models.ChangeBalanceRequest{
			ID:            walletID,
			Balance:       500,
			OperationType: models.Operation_type_deposit,
		}
		mockRepo.On("DepositUpdate", ctx, walletID, req.Balance).Return(nil)

		errResp := svc.ChangeWalletBalance(ctx, req)
		assert.Nil(t, errResp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("withdraw insufficient balance", func(t *testing.T) {
		req := models.ChangeBalanceRequest{
			ID:            walletID,
			Balance:       1000,
			OperationType: models.Operation_type_withdraw,
		}
		mockRepo.On("WithdrawUpdate", ctx, walletID, req.Balance).Return(repositories.ErrWalletNotEnoughBalance)

		errResp := svc.ChangeWalletBalance(ctx, req)
		assert.Equal(t, http.StatusBadRequest, errResp.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestWalletService_ChangeWalletBalance_Withdraw_NotFound(t *testing.T) {
	mockRepo := new(repositories.WalletRepoMock)
	svc := services.NewWalletService(mockRepo)
	ctx := context.Background()
	walletID := uuid.New()

	t.Run("withdraw wallet not found", func(t *testing.T) {
		req := models.ChangeBalanceRequest{
			ID:            walletID,
			Balance:       1000,
			OperationType: models.Operation_type_withdraw,
		}
		mockRepo.On("WithdrawUpdate", ctx, walletID, req.Balance).Return(repositories.ErrNoRowsForUpdate)

		errResp := svc.ChangeWalletBalance(ctx, req)
		assert.Equal(t, http.StatusNotFound, errResp.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestWalletService_ChangeWalletBalance_Withdraw_InternalError(t *testing.T) {
	mockRepo := new(repositories.WalletRepoMock)
	svc := services.NewWalletService(mockRepo)
	ctx := context.Background()
	walletID := uuid.New()

	t.Run("withdraw internal error", func(t *testing.T) {
		req := models.ChangeBalanceRequest{
			ID:            walletID,
			Balance:       1000,
			OperationType: models.Operation_type_withdraw,
		}
		mockRepo.On("WithdrawUpdate", ctx, walletID, req.Balance).Return(errors.New("db error"))

		errResp := svc.ChangeWalletBalance(ctx, req)
		assert.Equal(t, http.StatusInternalServerError, errResp.Code)
		mockRepo.AssertExpectations(t)
	})
}
