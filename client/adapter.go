package main

import (
	"context"
	"log"

	pb "kata_test/generated/protobuf/rates"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GetRatesService struct {
	rpc    pb.RatesServiceClient
	addr   string
	logger *zap.Logger
}

func NewUserService(addr string, logger *zap.Logger) *GetRatesService {

	serviceConfig := grpc.WithDefaultServiceConfig(`{
		"loadBalancingPolicy": "round_robin",
		"healthCheckConfig": {
		  "serviceName": ""
		}
	  }`)
	conn, err := grpc.NewClient(addr, serviceConfig, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error when creating grpc client: %v", err)
	}
	c := pb.NewRatesServiceClient(conn)

	return &GetRatesService{
		rpc:    c,
		addr:   addr,
		logger: logger,
	}

}

func (g *GetRatesService) GetRates(ctx context.Context, currency int) ([]Order, error) {

	rawRequest := &pb.RatesRequest{
		Currency: pb.Currency(currency),
	}

	rawResponse, err := g.rpc.GetRates(ctx, rawRequest)
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0, len(rawResponse.Orders))
	for _, rawOrder := range rawResponse.Orders {
		orders = append(orders, Order{
			Timestamp: rawOrder.Timestamp,
			Market:    rawOrder.Market,
			Ask:       rawOrder.Ask,
			Bid:       rawOrder.Bid,
		})
	}

	return orders, nil
}
