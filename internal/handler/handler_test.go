package handler_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"usdt/internal/handler"
	pb "usdt/internal/handler/pb"
	"usdt/internal/models"

	"github.com/stretchr/testify/assert"
)

type mockService struct {
	rate *models.Rate
	err  error
}

func (m *mockService) GetAndStoreRates(ctx context.Context) (*models.Rate, error) {
	return m.rate, m.err
}

func TestGetRates_Success(t *testing.T) {
	rate := &models.Rate{Ask: 10.0, Bid: 9.5, Timestamp: time.Date(2025, 7, 24, 0, 0, 0, 0, time.UTC)}
	svc := &mockService{rate: rate}
	server := &handler.Server{Service: svc}

	resp, err := server.GetRates(context.Background(), &pb.Empty{})
	assert.NoError(t, err)
	assert.Equal(t, rate.Ask, resp.Ask)
	assert.Equal(t, rate.Bid, resp.Bid)
	assert.Equal(t, "2025-07-24T00:00:00Z", resp.Timestamp)
}

func TestGetRates_Error(t *testing.T) {
	svc := &mockService{err: errors.New("some error")}
	server := &handler.Server{Service: svc}

	resp, err := server.GetRates(context.Background(), &pb.Empty{})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestHealthCheck(t *testing.T) {
	server := &handler.Server{}
	resp, err := server.HealthCheck(context.Background(), &pb.Empty{})
	assert.NoError(t, err)
	assert.Equal(t, "OK", resp.Status)
}
