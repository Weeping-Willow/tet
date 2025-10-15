package rates

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Weeping-Willow/tet/internal/utils"
	"github.com/pkg/errors"
)

type service struct {
	repo    Repository
	fetcher Fetcher
}

type fetcherResponse struct {
	CurrencyRate CurrencyRate
	Error        error
}

func NewService(repo Repository, fetcher Fetcher) Service {
	return &service{
		repo:    repo,
		fetcher: fetcher,
	}

}

func (s service) UpdateRates(ctx context.Context) error {
	utils.LoggerFromContext(ctx).Info("Updating rates")

	rates, errs, err := s.getLatestRates(ctx, preselectedCurrencies)
	if err != nil {
		return errors.Wrap(err, "get latest rates")
	}

	fmt.Println("Rates:", rates)
	fmt.Println("Errors:", errs)

	// Save rates to repository

	// save to repo update stats

	//TODO implement me
	panic("implement me")
}

func (s service) getLatestRates(ctx context.Context, currencyCodes []string) (rates []CurrencyRate, errs []string, err error) {
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
			errs = append(errs, resp.Error.Error())
			continue
		}

		rates = append(rates, resp.CurrencyRate)
	}

	logger.Info(fmt.Sprintf("Fetched %d currency rates, %d suceded, %d failed", len(currencyCodes), len(rates), len(errs)))

	if len(rates) == 0 && len(errs) > 0 {
		return nil, errs, errors.New("failed to fetch any currency rates")
	}

	return rates, errs, nil
}
