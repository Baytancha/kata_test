package rates

import (
	"context"
	pb "kata_test/generated/protobuf/rates"
	"kata_test/internal/domain"
	xlogger "kata_test/internal/infrastructure/logger"
	xtracer "kata_test/internal/infrastructure/tracer"
	service "kata_test/internal/service/rates"
	metrics "kata_test/metrics"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type RatesService interface {
	GetRates(ctx context.Context, currency domain.Currency) ([]domain.Order, error)
}

type RatesHandler struct {
	service RatesService
	logger  *zap.Logger
	metrics *metrics.PromMetrics
	tracer  trace.Tracer
	pb.UnimplementedRatesServiceServer
}

func NewHandler(srv *service.Rates) *RatesHandler {
	return &RatesHandler{
		service: srv,
		logger:  xlogger.Logger().Named("rates_handler"),
		metrics: metrics.Metrics(),
		tracer:  xtracer.Tracer(),
	}
}

func (h *RatesHandler) GetRates(ctx context.Context, req *pb.RatesRequest) (*pb.RatesResponse, error) {
	start := time.Now()

	defer func() {
		h.metrics.Response_duration.WithLabelValues("GetRates_handler_duration").
			Observe(time.Since(start).Seconds())
	}()

	h.metrics.ConsumeCounter.WithLabelValues("GetRates_handler_counter").Inc()

	ctx, span := h.tracer.Start(ctx, "Get_Rates_handler")

	defer span.End()

	var currency domain.Currency

	switch req.GetCurrency() {
	case 1:
		currency = "usdt"
	case 2:
		currency = "btc"
	default:
		currency = "usdt"
	}

	orders, err := h.service.GetRates(ctx, currency)
	if err != nil {
		h.logger.Error("error getting rates", zap.Error(err))
		return nil, err
	}

	results := &pb.RatesResponse{}
	for _, order := range orders {
		results.Orders = append(results.Orders, &pb.Order{
			Timestamp: order.Timestamp,
			Market:    order.Market,
			Ask:       order.Ask,
			Bid:       order.Bid,
		})
	}

	return results, nil
}
