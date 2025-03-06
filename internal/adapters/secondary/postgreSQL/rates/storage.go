package storage

import (
	"context" // New import
	"database/sql"
	"errors"
	"kata_test/internal/domain"
	xlogger "kata_test/internal/infrastructure/logger"
	xtracer "kata_test/internal/infrastructure/tracer"
	metrics "kata_test/metrics"
	"time"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type OrderModel struct {
	logger  *zap.Logger
	metrics *metrics.PromMetrics
	tracer  trace.Tracer
	DB      *sqlx.DB
}

func NewOrderModel(db *sqlx.DB) *OrderModel {
	return &OrderModel{
		DB:      db,
		logger:  xlogger.Logger().Named("rates_storage"),
		metrics: metrics.Metrics(),
		tracer:  xtracer.Tracer(),
	}
}

func (m OrderModel) SaveOrder(ctx context.Context, order domain.Order) error {
	start := time.Now()
	defer func() {
		m.metrics.DB_duration.WithLabelValues("GetRates_db_duration").
			Observe(time.Since(start).Seconds())
	}()

	ctx, span := m.tracer.Start(ctx, "Get_Rates_db")

	defer span.End()

	tx, err := m.DB.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Join(domain.ErrTx, err)
	}

	defer func() {
		rollbackErr := tx.Rollback()
		if !errors.Is(rollbackErr, sql.ErrTxDone) {
			m.logger.Info("rollback error", zap.Error(rollbackErr))
		}
	}()

	query := `
	INSERT INTO rates (timestamp, market, ask, bid)
	 VALUES ($1, $2, $3, $4)
	 RETURNING timestamp
`

	args := []interface{}{
		time.Unix(order.Timestamp, 0),
		order.Market,
		order.Ask,
		order.Bid,
	}
	var timestamp time.Time
	err = tx.QueryRowxContext(ctx, query, args...).Scan(&timestamp)

	if err != nil {
		return errors.Join(domain.ErrRecordNotFound, err)
	}
	order.Timestamp = timestamp.Unix()
	if err = tx.Commit(); err != nil {
		return errors.Join(domain.ErrTx, err)
	}

	return nil
}
