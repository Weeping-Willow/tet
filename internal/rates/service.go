package rates

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Weeping-Willow/tet/internal/objects"
	"github.com/Weeping-Willow/tet/internal/utils"
	"github.com/pkg/errors"
)

type service struct {
	repo    Repository
	fetcher Fetcher
}

type fetcherResponse struct {
	CurrencyRate objects.CurrencyRate
	Error        error
}

func NewService(repo Repository, fetcher Fetcher) Service {
	return &service{
		repo:    repo,
		fetcher: fetcher,
	}

}

func (s service) UpdateRates(ctx context.Context) error {
	logger := utils.LoggerFromContext(ctx)
	logger.Info("Updating rates")

	rates, _, err := s.getLatestRates(ctx, preselectedCurrencies)
	if err != nil {
		return errors.Wrap(err, "get latest rates")
	}

	logger.Info(fmt.Sprintf("Saving rates"))

	err = s.repo.UpdateRates(ctx, rates)
	if err != nil {
		return errors.Wrap(err, "update rates in repo")
	}

	logger.Info("Rates saved successfully")

	return nil
}

func (s service) getLatestRates(ctx context.Context, currencyCodes []string) (rates []objects.CurrencyRate, individualFetchingErrors []string, err error) {
	if len(currencyCodes) == 0 {
		return nil, nil, errors.New("no currency codes provided")
	}

	logger := utils.LoggerFromContext(ctx)

	logger.Info("Fetching latest currency rates")

	wg := sync.WaitGroup{}
	ch := make(chan fetcherResponse, len(currencyCodes))

	for _, code := range currencyCodes {
		wg.Add(1)
		go func(code string) {
			defer wg.Done()

			logger.Info("Fetching latest currency rate", "code", code)
			tn := time.Now()
			res, err := s.fetcher.GetCurrencyRate(ctx, code)

			select {
			case ch <- fetcherResponse{
				CurrencyRate: res,
				Error:        errors.Wrapf(err, "fetch currency rate for %s", code),
			}:
			case <-ctx.Done():
				logger.Warn("Context cancelled, stopping fetch", "code", code)

				return
			}

			if err == nil {
				logger.Info("Fetched latest currency rate", "code", code, "duration", time.Since(tn).String())
			}

			if err != nil {
				logger.Error("Fetch latest currency rate", "code", code, "error", err.Error())
			}
		}(code)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for resp := range ch {
		if resp.Error != nil {
			individualFetchingErrors = append(individualFetchingErrors, resp.Error.Error())
			continue
		}

		rates = append(rates, resp.CurrencyRate)
	}

	logger.Info(fmt.Sprintf("Fetched %d currency rates, %d suceded, %d failed", len(currencyCodes), len(rates), len(individualFetchingErrors)))

	if len(rates) == 0 && len(individualFetchingErrors) > 0 {
		return nil, individualFetchingErrors, errors.New("failed to fetch any currency rates")
	}

	return rates, individualFetchingErrors, nil
}
