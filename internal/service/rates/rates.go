package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kata_test/internal/domain"
	xlogger "kata_test/internal/infrastructure/logger"
	"kata_test/metrics"
	"net/http"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Storage interface {
	SaveOrder(ctx context.Context, order domain.Order) error
}

type Rates struct {
	client  *http.Client
	storage Storage
	logger  *zap.Logger
	metrics *metrics.PromMetrics
}

func NewRates(storage Storage) *Rates {
	return &Rates{
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   5 * time.Second,
		},
		storage: storage,
		logger:  xlogger.Logger(),
	}
}

var (
	garantexAPIURL = "https://garantex.org/api/v2/depth"
	currencies     map[domain.Currency][]domain.Market
)

type OrdersDTO struct {
	Timestamp int64 `json:"timestamp"`
	Asks      []Ask `json:"asks"`
	Bids      []Ask `json:"bids"`
}

type Ask struct {
	Price float64 `json:"price,string"`
}

type Bid struct {
	Price float64 `json:"price,string"`
}

func init() {
	currencies = make(map[domain.Currency][]domain.Market)
	currencies["usdt"] = []domain.Market{
		"usdta7a5",
		"usdtkgs",
		"usdtrub",
		"usdteur",
		"usdtusd",
	}
}

func (svc *Rates) GetRates(ctx context.Context, currency domain.Currency) ([]domain.Order, error) {
	var orders []domain.Order
	var markets []domain.Market

	switch currency {
	case "usdt":
		markets = currencies[currency]
		orders = make([]domain.Order, len(markets))
	default:
		return nil, errors.New("invalid currency type")
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(markets))
	for i, name := range markets {

		url, err := url.Parse(garantexAPIURL)
		if err != nil {
			return nil, err
		}
		values := url.Query()
		values.Add("market", string(name))
		url.RawQuery = values.Encode()

		req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)

		if err != nil {
			return nil, err
		}

		res, err := svc.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		data := OrdersDTO{}

		err = json.Unmarshal(body, &data)
		if err != nil {
			svc.logger.Error(fmt.Sprintf("error unmarshaling data: %v", err))
			return nil, err
		}

		order := domain.Order{}
		order.Market = string(name)
		order.Timestamp = data.Timestamp
		order.Ask = data.Asks[0].Price
		order.Bid = data.Bids[len(data.Bids)-1].Price

		go func(i int, order domain.Order) {
			defer wg.Done()
			err = svc.storage.SaveOrder(ctx, order)
			if err != nil {
				svc.logger.Info("error saving order", zap.Error(err))
			}
			orders[i] = order
		}(i, order)

	}
	wg.Wait()
	return orders, nil
}
