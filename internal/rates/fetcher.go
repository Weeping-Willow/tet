package rates

import (
	"context"
	"net/http"
	"time"

	"github.com/Weeping-Willow/tet/internal/config"
)

type ecbRSSFetcher struct {
	url    string
	client *http.Client
}

func NewEcbRssFetcher(cfg config.Config) Fetcher {
	return &ecbRSSFetcher{
		url: cfg.ExternalServices.EcbRssURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (e ecbRSSFetcher) GetCurrencyRate(ctx context.Context, currencyCode string) (CurrencyRate, error) {
	time.Sleep(time.Second)

	return CurrencyRate{
		Currency: currencyCode,
		Rate:     1.2345,
		Date:     time.Now().Format("2006-01-02"),
	}, nil
}
