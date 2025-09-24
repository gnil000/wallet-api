package migrations

import (
	"embed"
	"fmt"
	"wallet-api/config"
	"wallet-api/pkg/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	glog "go.finelli.dev/gooseloggers/zerolog"
)

//go:embed *.sql
var embedMigrations embed.FS

type Migrator struct {
	logger zerolog.Logger
	cfg    config.DB
}

func NewMigrator(logger zerolog.Logger, cfg config.DB) *Migrator {
	return &Migrator{logger: logger, cfg: cfg}
}

func (m *Migrator) Migrate() {
	err := m.migrate()
	if err != nil {
		panic(fmt.Errorf("failed to migrate: %w", err))
	}
}

func (m *Migrator) migrate() error {
	db, err := goose.OpenDBWithDriver("pgx", m.cfg.WalletDB.ConnectionString)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}

	goose.SetTableName("vehicles_api_migrations")
	goose.SetBaseFS(embedMigrations)
	log := logger.WithModule(m.logger, "migrations")
	goose.SetLogger(glog.GooseZerologLogger(&log))

	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("failed to make migrations: %w", err)
	}

	return nil
}
