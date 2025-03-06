package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kata_test/internal/domain"
	"kata_test/metrics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SaveOrder(ctx context.Context, order domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func TestGetRates(t *testing.T) {
	var setupMocks func(*MockStorage)

	setupMocks = func(storage *MockStorage) {
		storage.On("SaveOrder", mock.Anything, mock.Anything).
			Return(nil)

	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a response from the Garantex API
		response := OrdersDTO{
			Timestamp: 999,
			Asks: []Ask{
				{Price: 1.0},
				{Price: 1.1},
				{Price: 1.2},
			},
			Bids: []Ask{
				{Price: 0.9},
				{Price: 0.8},
				{Price: 0.7},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	storage := &MockStorage{}
	setupMocks(storage)

	rates := &Rates{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		storage: storage,
		logger:  zap.NewNop(),
		metrics: metrics.Metrics(),
	}

	oldURL := garantexAPIURL
	defer func() {
		garantexAPIURL = oldURL
	}()

	garantexAPIURL = server.URL

	currency := domain.Currency("usdt")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orders, err := rates.GetRates(ctx, currency)
	if err != nil {
		t.Fatalf("GetRates returned an error: %v", err)
	}

	assert.Equal(t, len(currencies[currency]), len(orders))

	if !storage.AssertExpectations(t) {
		t.Error("Expected SaveOrder to be called")
	}

	expectedOrder := domain.Order{
		Timestamp: 999,
		Ask:       1.0,
		Bid:       0.7,
	}
	for _, order := range orders {
		assert.Equal(t, expectedOrder.Timestamp, order.Timestamp)
		assert.Equal(t, expectedOrder.Ask, order.Ask)
		assert.Equal(t, expectedOrder.Bid, order.Bid)
	}

}

func TestGetRatesApi(t *testing.T) {
	var setupMocks func(*MockStorage)

	setupMocks = func(storage *MockStorage) {
		storage.On("SaveOrder", mock.Anything, mock.Anything).
			Return(nil)

	}

	storage := &MockStorage{}
	setupMocks(storage)

	rates := &Rates{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		storage: storage,
		logger:  zap.NewNop(),
		metrics: metrics.Metrics(),
	}

	currency := domain.Currency("usdt")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orders, err := rates.GetRates(ctx, currency)
	if err != nil {
		t.Fatalf("GetRates returned an error: %v", err)
	}

	for _, order := range orders {
		fmt.Println(order)
	}

}
