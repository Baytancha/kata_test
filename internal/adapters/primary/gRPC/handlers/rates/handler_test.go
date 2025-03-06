package rates

import (
	"context"
	"fmt"
	pb "kata_test/generated/protobuf/rates"
	"kata_test/internal/domain"
	metrics "kata_test/metrics"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.uber.org/zap"
)

type MockRates struct {
	mock.Mock
}

func (m *MockRates) GetRates(ctx context.Context, currency domain.Currency) ([]domain.Order, error) {
	args := m.Called(ctx, currency)
	return args.Get(0).([]domain.Order),
		args.Error(1)
}

func TestGetRatesHandler(t *testing.T) {
	var setupMocks func(*MockRates)

	sr := tracetest.NewSpanRecorder()

	exp := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(sr),
		trace.WithBatcher(exp),
	)
	defer tp.Shutdown(context.Background())

	testRequest := &pb.RatesRequest{
		Currency: pb.Currency(1),
	}

	t.Run("happy path", func(t *testing.T) {
		exp.Reset()

		setupMocks = func(service *MockRates) {
			service.On("GetRates",
				mock.Anything,
				domain.Currency("usdt")).
				Return([]domain.Order{
					{Market: "usdt",
						Timestamp: 1},
					{Market: "btc",
						Timestamp: 2},
					{Market: "eth",
						Timestamp: 3}}, nil)

		}
		mockService := &MockRates{}
		setupMocks(mockService)

		handler := &RatesHandler{
			service: mockService,
			logger:  zap.NewNop(),
			metrics: metrics.Metrics(),
			tracer:  tp.Tracer("test"),
		}

		resp, err := handler.GetRates(context.Background(), testRequest)

		if !mockService.AssertExpectations(t) {
			t.Error("Expected GetRates to be called")
		}
		assert.NoError(t, err)
		assert.NotNil(t, resp)

		count := testutil.ToFloat64(handler.metrics.ConsumeCounter)
		assert.Equal(t, 1.0, count)

		exp.ExportSpans(context.Background(), sr.Ended())
		spans := exp.GetSpans()
		assert.Greater(t, len(spans), 0, "No spans recorded")

		span := spans[0]
		assert.Equal(t, "Get_Rates_handler", span.Name)

	})

	t.Run("error case", func(t *testing.T) {
		exp.Reset()

		setupMocks = func(service *MockRates) {
			service.On("GetRates",
				mock.Anything,
				domain.Currency("usdt")).
				Return([]domain.Order{}, fmt.Errorf("error"))

		}

		mockService := &MockRates{}
		setupMocks(mockService)

		handler := &RatesHandler{
			service: mockService,
			logger:  zap.NewNop(),
			metrics: metrics.Metrics(),
			tracer:  tp.Tracer("test"),
		}
		resp, err := handler.GetRates(context.Background(), testRequest)

		if !mockService.AssertExpectations(t) {
			t.Error("Expected GetRates to be called")
		}

		assert.Error(t, err)
		assert.Nil(t, resp)

		count := testutil.ToFloat64(handler.metrics.ConsumeCounter)
		assert.Equal(t, 2.0, count)

		exp.ExportSpans(context.Background(), sr.Ended())
		spans := exp.GetSpans()
		assert.Greater(t, len(spans), 0, "No spans recorded")

		span := spans[0]
		assert.Equal(t, "Get_Rates_handler", span.Name)

	})

}
