package service_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"usdt/internal/models"
	"usdt/internal/service"

	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	rate *models.Rate
	err  error
}

func (m *mockClient) GetDepth() (*models.Rate, error) {
	return m.rate, m.err
}

type mockRepo struct {
	err error
}

func (m *mockRepo) SaveRate(ctx context.Context, rate *models.Rate) error {
	return m.err
}

func TestGetAndStoreRates_Success(t *testing.T) {
	mockRate := &models.Rate{Ask: 10.0, Bid: 9.5, Timestamp: time.Now()}
	c := &mockClient{rate: mockRate}
	r := &mockRepo{}
	svc := &service.RateService{Client: c, Repo: r}

	result, err := svc.GetAndStoreRates(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, mockRate, result)
}

func TestGetAndStoreRates_ClientError(t *testing.T) {
	c := &mockClient{err: errors.New("client error")}
	r := &mockRepo{}
	svc := &service.RateService{Client: c, Repo: r}

	result, err := svc.GetAndStoreRates(context.Background())
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetAndStoreRates_RepoError(t *testing.T) {
	mockRate := &models.Rate{Ask: 10.0, Bid: 9.5, Timestamp: time.Now()}
	c := &mockClient{rate: mockRate}
	r := &mockRepo{err: errors.New("repo error")}
	svc := &service.RateService{Client: c, Repo: r}

	result, err := svc.GetAndStoreRates(context.Background())
	assert.Error(t, err)
	assert.Nil(t, result)
}
