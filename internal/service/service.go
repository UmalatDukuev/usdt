package service

import (
	"context"
	"usdt/internal/models"
)

type RateServiceInterface interface {
	GetAndStoreRates(ctx context.Context) (*models.Rate, error)
}

type GrinexClientInterface interface {
	GetDepth() (*models.Rate, error)
}

type RepoInterface interface {
	SaveRate(ctx context.Context, rate *models.Rate) error
}

type RateService struct {
	Client GrinexClientInterface
	Repo   RepoInterface
}

// GetAndStoreRates fetches the current rate from the client and saves it in the repository.
func (s *RateService) GetAndStoreRates(ctx context.Context) (*models.Rate, error) {
	rate, err := s.Client.GetDepth()
	if err != nil {
		return nil, err
	}

	if err := s.Repo.SaveRate(ctx, rate); err != nil {
		return nil, err
	}

	return rate, nil
}
