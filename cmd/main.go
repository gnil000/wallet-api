package main

import (
	"context"
	"wallet-api/config"
	"wallet-api/pkg/database"
	"wallet-api/pkg/httpserver"
	"wallet-api/pkg/logger"
	"wallet-api/src/database/migrations"
	"wallet-api/src/database/repositories"
	"wallet-api/src/handlers"
	"wallet-api/src/services"

	"github.com/rs/zerolog"
)

type appKeyType struct{}

var appKey = appKeyType{}

func main() {
	cfg := config.GetConfig()

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), appKey, cfg.App))
	defer cancel()

	connPool := database.NewConnectionPool(ctx, cfg.Database.WalletDB)

	log := logger.NewLogger(zerolog.InfoLevel)
	migrator := migrations.NewMigrator(log, cfg.Database)

	migrator.Migrate()

	//---

	wallet_repo := repositories.NewWalletRepo(connPool, log)
	wallet_service := services.NewWalletService(wallet_repo)
	wallet_handler := handlers.NewWalletHandler(wallet_service)

	server := httpserver.NewServer(log, cfg.Server)
	wallet_handler.Register(server)

	server.Serve(ctx)

}
