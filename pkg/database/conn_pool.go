package database

import (
	"context"
	"errors"
	"log"
	"wallet-api/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionPool interface {
	GetConnection(ctx context.Context) (Connection, error)
	Close()
}

type Connection interface {
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Release()
}

type pool struct {
	*pgxpool.Pool
}

func NewConnectionPool(ctx context.Context, c ConnectionConfig) ConnectionPool {
	cfg, err := pgxpool.ParseConfig(c.ConnectionString)
	if err != nil {
		log.Fatal("Failed to parse connection string", err)
	}
	cfg.ConnConfig.ConnectTimeout = c.ConnectionTimeout
	if app, ok := utils.ContextApp(ctx); ok {
		cfg.ConnConfig.RuntimeParams["application_name"] = app
	}
	if c.UsePGBouncer {
		cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	}
	if c.PoolMinConns > 0 {
		cfg.MinConns = c.PoolMinConns
	}
	if c.PoolMaxConns > 0 {
		cfg.MaxConns = c.PoolMaxConns
	}

	p, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create connection pool", err)
	}
	pool := &pool{p}
	return pool
}

func (p *pool) GetConnection(ctx context.Context) (Connection, error) {
	conn, err := p.Pool.Acquire(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		//TODO:
	}
	return conn, err
}

func (p *pool) Close() {
	//db, host := p.Config().ConnConfig.Database, p.Config().ConnConfig.Host
	//log.Info("Closing connection pool", "db", db, "host", host)
	p.Pool.Close()
}
