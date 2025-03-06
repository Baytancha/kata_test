package main

import (
	"context"
	"flag"
	"kata_test/config"
	xlogger "kata_test/internal/infrastructure/logger"
	"kata_test/metrics"
	"log"

	_ "net/http/pprof"

	_ "github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	tp := config.TracerInit()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("error shutting down tracer provider: %v", err)
		}
	}()

	xlogger.BuildLogger(cfg.LogLevel)
	logger := xlogger.Logger().Named("main")

	flag.StringVar(&cfg.Db.Dsn, "db-dsn", cfg.Db.Dsn, "PostgreSQL DSN")
	flag.IntVar(&cfg.Db.MaxOpenConns, "db-max-open-conns", cfg.Db.MaxOpenConns, "PostgreSQL max open connections")
	flag.IntVar(&cfg.Db.MaxIdleConns, "db-max-idle-conns", cfg.Db.MaxIdleConns, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.Db.MaxIdleTime, "db-max-idle-time", cfg.Db.MaxIdleTime, "PostgreSQL max connection idle time")

	app := NewApp(cfg, logger)
	app.Run()

	metrics.BuildMetrics()

	go func() {
		err := metrics.MetricsServer(metrics.Metrics()).ListenAndServe()
		logger.Fatal("error starting metrics server", zap.Error(err))
	}()

	logger.Info("server started")
	err = app.Serve()
	if err != nil {
		logger.Fatal("error starting server", zap.Error(err))
	}
	logger.Info("server stopped")
}
