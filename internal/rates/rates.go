package rates

import "context"

type Repository interface {
}

type Fetcher interface {
	//GetCurrencyRate(ctx context.Context, currencyCode string) (float64, error)
}

type Service interface {
	UpdateRates(ctx context.Context) error
	//GetLatestRates(ctx context.Context) (any, error)
	//GetCurrencyHistory(ctx context.Context, currency string) (any, error)
}
