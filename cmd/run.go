package main

import (
	"kata_test/config"
	rateHandler "kata_test/internal/adapters/primary/gRPC/handlers/rates"
	server "kata_test/internal/adapters/primary/gRPC/server"
	storage "kata_test/internal/adapters/secondary/postgreSQL/rates"
	"kata_test/internal/db"
	service "kata_test/internal/service/rates"

	"go.uber.org/zap"
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
	server *server.GrpcServer
}

func NewApp(conf *config.Config, logger *zap.Logger) *App {
	return &App{
		cfg:    conf,
		logger: logger,
	}
}

func (app *App) Serve() error {

	app.server.Shutdown()

	err := app.server.SetupAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) Run() {
	DbCfg := &db.DBconfig{
		Dsn:          a.cfg.Db.Dsn,
		MaxOpenConns: a.cfg.Db.MaxOpenConns,
		MaxIdleConns: a.cfg.Db.MaxIdleConns,
		MaxIdleTime:  a.cfg.Db.MaxIdleTime,
	}
	rpcCfg := &server.GRPConfig{
		Port: a.cfg.Grpc.Port,
		Addr: a.cfg.Grpc.Address,
	}

	dbx, err := db.NewSqlDB(DbCfg, a.logger)
	if err != nil {
		a.logger.Fatal("error init db", zap.Error(err))
	}

	rateStorage := storage.NewOrderModel(dbx)
	rateService := service.NewRates(rateStorage)
	rateHandler := rateHandler.NewHandler(rateService)
	a.server = server.New(
		rpcCfg,
		&server.Handlers{
			Rates: rateHandler})

}
