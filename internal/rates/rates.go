package rates

import "context"

type Repository interface {
}

type Fetcher interface {
	GetCurrencyRate(ctx context.Context, currencyCode string) (CurrencyRate, error)
}

type Service interface {
	UpdateRates(ctx context.Context) error
	//GetLatestRates(ctx context.Context) (any, error)
	//GetCurrencyHistory(ctx context.Context, currency string) (any, error)
}

type CurrencyRate struct {
	Rate     float64 `json:"rate"`
	Currency string  `json:"currency"`
	Date     string  `json:"date"`
}

var preselectedCurrencies = []string{
	"USD",
	"EUR",
	"GBP",
	"JPY",
	"AUD",
	"CAD",
	"CNY",
	"CHF",
	"SEK",
	"THB",
}
