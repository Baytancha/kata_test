package main

import (
	"context"
	"fmt"
	"kata_test/config"
	xlogger "kata_test/internal/infrastructure/logger"
	"log"
	"time"

	"go.uber.org/zap"
	_ "google.golang.org/grpc/health"
)

var addr string

func main() {

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	if cfg.Grpc.ClientAddr == ":8181" {
		addr = ":8181"
	} else {
		addr = cfg.Grpc.ClientAddr
	}

	xlogger.BuildLogger(cfg.LogLevel)
	logger := xlogger.Logger().Named("main")
	logger.Info("server started")
	client := NewUserService(addr, logger)

	for {
		time.Sleep(time.Second)
		fmt.Println("Enter currency: ")
		var currency int

		resp, err := client.GetRates(context.Background(), currency)
		if err != nil {
			logger.Fatal("error get rates", zap.Error(err))
		}
		fmt.Printf("%+v\n", resp)
	}

}
