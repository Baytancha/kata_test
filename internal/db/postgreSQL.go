package db

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DBconfig struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func NewSqlDB(cfg *DBconfig, logger *zap.Logger) (*sqlx.DB, error) {
	connConfig, _ := pgx.ParseConfig(cfg.Dsn)
	pgxdb := stdlib.OpenDB(*connConfig)

	db := sqlx.NewDb(pgxdb, "pgx")
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	if err := db.Ping(); err != nil {
		logger.Error("error starting db", zap.Error(err))
		return nil, err
	}

	return db, nil

}
