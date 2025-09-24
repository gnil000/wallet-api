package config

import (
	_ "embed"
	"errors"
	"log"
	"os"
	"strings"
	"wallet-api/pkg/database"
	"wallet-api/pkg/httpserver"

	"gopkg.in/yaml.v3"
)

const (
	envWalletDB = "PG_WALLET"
)

type Config struct {
	App      string                  `yaml:"app"`
	Stack    string                  `yaml:"stack"`
	Database DB                      `yaml:"db"`
	Server   httpserver.ServerConfig `yaml:"server"`
}

type DB struct {
	WalletDB database.ConnectionConfig `yaml:"walletDB"`
}

var (
	//go:embed config.yaml
	file    string
	decoder = yaml.NewDecoder(strings.NewReader(file))
)

func GetConfig() Config {
	var config Config
	if err := decoder.Decode(&config); err != nil {
		log.Fatal("Config decoding error", err)
	}
	if envVal, exist := os.LookupEnv(envWalletDB); exist {
		log.Println("Found database in env")
		config.Database.WalletDB.ConnectionString = envVal
	}
	if err := validate(config); err != nil {
		log.Fatal("Invalid configuration", err)
	}
	return config
}

func validate(cfg Config) error {
	switch {
	case cfg.Database.WalletDB.ConnectionString == "":
		return errors.New("wallet db connection string is empty")
	case cfg.Server.Port == 0:
		return errors.New("server port is empty")
	default:
		return nil
	}
}
