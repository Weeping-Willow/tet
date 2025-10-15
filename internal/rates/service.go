package rates

import (
	"context"

	"github.com/Weeping-Willow/tet/internal/utils"
)

type service struct {
	repo    Repository
	fetcher Fetcher
}

func NewService(repo Repository, fetcher Fetcher) Service {
	return &service{
		repo:    repo,
		fetcher: fetcher,
	}

}

func (s service) UpdateRates(ctx context.Context) error {
	utils.LoggerFromContext(ctx).Info("Updating rates")
	//TODO implement me
	panic("implement me")
}
