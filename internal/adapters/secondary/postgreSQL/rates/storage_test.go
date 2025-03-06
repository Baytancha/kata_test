package storage

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	"kata_test/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestMetricsModel_success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectQuery(`
	INSERT INTO rates (timestamp, market, ask, bid)
	 VALUES ($1, $2, $3, $4)
	 RETURNING timestamp
`).WithArgs(
		time.Unix(0, 0),
		"usdtrub",
		21.34,
		21.36,
	).WillReturnRows(sqlmock.NewRows([]string{"timestamp"}).AddRow(time.Now()))

	mock.ExpectCommit()

	dbx := sqlx.NewDb(db, "sqlmock")
	metricsModel := NewOrderModel(dbx)

	err = metricsModel.SaveOrder(context.Background(), domain.Order{
		Timestamp: 0,
		Market:    "usdtrub",
		Ask:       21.34,
		Bid:       21.36,
	})

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
