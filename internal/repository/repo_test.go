package repository_test

import (
	"context"
	"testing"
	"time"
	"usdt/internal/models"
	"usdt/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"go.opentelemetry.io/otel"
)

func TestSaveRate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	r := &repository.Repo{
		DB:     db,
		Tracer: otel.Tracer("test-tracer"),
	}

	mock.ExpectExec("^INSERT INTO rates").
		WithArgs(10.0, 9.5, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = r.SaveRate(context.Background(), &models.Rate{Ask: 10.0, Bid: 9.5, Timestamp: time.Now()})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
