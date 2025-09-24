package services

import (
	"context"
	"errors"
	"net/http"
	"wallet-api/src/database/repositories"
	"wallet-api/src/models"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type WalletService interface {
	GetWalletByID(ctx context.Context, id uuid.UUID) (models.GetBalanceResponse, *models.ErrorResponse)
	ChangeWalletBalance(ctx context.Context, changeBalanceReq models.ChangeBalanceRequest) *models.ErrorResponse
	GetWallets(ctx context.Context) ([]models.GetWalletsResponse, *models.ErrorResponse)
}

type walletService struct {
	walletRepo repositories.WalletRepo
}

func NewWalletService(walletRepo repositories.WalletRepo) WalletService {
	return &walletService{walletRepo: walletRepo}
}

func (s *walletService) GetWalletByID(ctx context.Context, id uuid.UUID) (models.GetBalanceResponse, *models.ErrorResponse) {
	var wallet models.GetBalanceResponse
	walletEntity, err := s.walletRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrWalletNotFound) {
			return models.GetBalanceResponse{}, &models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "wallet not found",
			}
		}
		return models.GetBalanceResponse{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	if err = copier.Copy(&wallet, &walletEntity); err != nil {
		return models.GetBalanceResponse{}, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return wallet, nil
}

func (s *walletService) ChangeWalletBalance(ctx context.Context, changeBalanceReq models.ChangeBalanceRequest) *models.ErrorResponse {
	var err error
	switch changeBalanceReq.OperationType {
	case models.Operation_type_deposit:
		err = s.walletRepo.DepositUpdate(ctx, changeBalanceReq.ID, changeBalanceReq.Balance)
	case models.Operation_type_withdraw:
		err = s.walletRepo.WithdrawUpdate(ctx, changeBalanceReq.ID, changeBalanceReq.Balance)
	}
	if err != nil {
		if errors.Is(err, repositories.ErrWalletNotEnoughBalance) {
			return &models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "not enough balance",
			}
		}
		if errors.Is(err, repositories.ErrNoRowsForUpdate) {
			return &models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "wallet not found",
			}
		}
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return nil
}

func (s *walletService) GetWallets(ctx context.Context) ([]models.GetWalletsResponse, *models.ErrorResponse) {
	var wallets []models.GetWalletsResponse
	walletEntities, err := s.walletRepo.GetWallets(ctx)
	if err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	if err = copier.Copy(&wallets, &walletEntities); err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return wallets, nil
}
