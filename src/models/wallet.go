package models

import "github.com/google/uuid"

const (
	Operation_type_deposit  = "DEPOSIT"
	Operation_type_withdraw = "WITHDRAW"
)

type Wallet struct {
	ID      uuid.UUID `json:"id"`
	Balance int64     `json:"balance"`
}

type GetWalletsResponse struct {
	ID uuid.UUID `json:"id"`
}

type GetBalanceResponse struct {
	Balance int64 `json:"balance"`
}

type ChangeBalanceRequest struct {
	ID            uuid.UUID `json:"valletId"`
	Balance       int64     `json:"amount"`
	OperationType string    `json:"operationType"`
}
