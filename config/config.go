package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Db struct {
		Dsn          string `env:"DB_DSN" envDefault:"postgres://postgres:password@localhost:5432/kata_test?sslmode=disable" `
		MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"100" `
		MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
		MaxIdleTime  string `env:"DB_MAX_IDLE_TIME" envDefault:"5m"`
	}

	Grpc struct {
		Address    string `env:"gRPC_ADDRESS" envDefault:":8181"`
		ClientAddr string `env:"gRPC_CLIENT_ADDRESS" envDefault:":8181" `
		Port       int    `env:"gRPC_PORT" envDefault:"8181"`
	}

	LogLevel          string `env:"LOG_LEVEL" envDefault:"DEBUG"`
	EnableDebugServer bool   `env:"ENABLE_DEBUG_SERVER" envDefault:"true"`
	DebugServerAddr   string `env:"DEBUG_SERVER_ADDR" envDefault:":6000"`
}

func ReadConfig() (*Config, error) {
	config := Config{}

	err := env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	return &config, err
}
