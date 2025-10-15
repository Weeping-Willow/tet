package rates

import (
	"context"

	"github.com/Weeping-Willow/tet/internal/objects"
)

type Repository interface {
	UpdateRates(ctx context.Context, rates []objects.CurrencyRate) error
}

type Fetcher interface {
	GetCurrencyRate(ctx context.Context, currencyCode string) (objects.CurrencyRate, error)
}

type Service interface {
	UpdateRates(ctx context.Context) error
	//GetLatestRates(ctx context.Context) (any, error)
	//GetCurrencyHistory(ctx context.Context, currency string) (any, error)
}

var preselectedCurrencies = []string{
	"USD",
	"PLN",
	"GBP",
	"JPY",
	"AUD",
	"CAD",
	"CNY",
	"CHF",
	"SEK",
	"THB",
}
