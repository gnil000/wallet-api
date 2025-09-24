package repositories

import (
	"context"
	"errors"
	"strings"
	"time"
	"wallet-api/pkg/database"
	"wallet-api/pkg/logger"
	"wallet-api/src/database/entities"
	"wallet-api/src/database/queries"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
)

const (
	maxRetries       = 20
	notEnoughBalance = "not enough balance"
)

var (
	module                    = "repo_wallet"
	ErrNoRowsForUpdate        = errors.New("row not found for affect")
	ErrWalletNotFound         = errors.New("wallet not found")
	ErrWalletNotEnoughBalance = errors.New("not enough balance")
)

type WalletRepo interface {
	FindByID(ctx context.Context, id uuid.UUID) (entities.Wallet, error)
	WithdrawUpdate(ctx context.Context, id uuid.UUID, amount int64) error
	DepositUpdate(ctx context.Context, id uuid.UUID, amount int64) error
	GetWallets(ctx context.Context) ([]entities.Wallet, error)
}

type walletRepository struct {
	pool database.ConnectionPool
	log  zerolog.Logger
}

func NewWalletRepo(pool database.ConnectionPool, log zerolog.Logger) WalletRepo {
	return &walletRepository{pool: pool, log: logger.WithModule(log, module)}
}

func (r *walletRepository) FindByID(ctx context.Context, id uuid.UUID) (entities.Wallet, error) {
	var err error
	connection, err := r.pool.GetConnection(ctx)
	if err != nil {
		return entities.Wallet{}, err
	}
	defer connection.Release()
	var wallet entities.Wallet
	err = connection.QueryRow(ctx, queries.FindWallet, id).Scan(
		&wallet.Balance,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Wallet{}, ErrWalletNotFound
		}
		return entities.Wallet{}, err
	}
	return wallet, nil
}

func (r *walletRepository) GetWallets(ctx context.Context) ([]entities.Wallet, error) {
	var err error
	connection, err := r.pool.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer connection.Release()
	var wallets []entities.Wallet
	rows, err := connection.Query(ctx, queries.GetWallets)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var wallet entities.Wallet
		err := rows.Scan(
			&wallet.ID,
		)
		if err != nil {
			continue
		}
		wallets = append(wallets, wallet)
	}
	return wallets, nil
}

func (r *walletRepository) DepositUpdate(ctx context.Context, id uuid.UUID, amount int64) error {
	var attempts int = 0
	var err error
	for {
		r.log.Debug().Msg("operation start deposit balance")
		if err = r.changeBalanceTx(ctx, id, amount, queries.UpdateDepositWallet); err == nil {
			r.log.Debug().Msg("operation deposit balance success")
			return nil
		}
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "40001" {
			if attempts < maxRetries {
				time.Sleep(time.Duration(attempts*100) * time.Millisecond)
				attempts++
				continue
			}
		}
		return err
	}
}

func (r *walletRepository) WithdrawUpdate(ctx context.Context, id uuid.UUID, amount int64) error {
	var attempts int = 0
	var err error
	for {
		r.log.Debug().Msg("operation start withdraw balance")
		if err = r.changeBalanceTx(ctx, id, amount, queries.UpdateWithdrawWallet); err == nil {
			r.log.Debug().Msg("operation withdraw balance success")
			return nil
		}
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "40001" {
			if attempts < maxRetries {
				time.Sleep(time.Duration(attempts*100) * time.Millisecond)
				attempts++
				continue
			}
		}
		return err
	}
}

func (r *walletRepository) changeBalanceTx(ctx context.Context, id uuid.UUID, amount int64, query string) error {
	var err error
	connection, err := r.pool.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()
	tx, err := connection.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	tag, err := connection.Exec(ctx, query, id, amount)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "P0002" {
			return ErrNoRowsForUpdate
		}
		if strings.Contains(err.Error(), notEnoughBalance) {
			return ErrWalletNotEnoughBalance
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoRowsForUpdate
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
