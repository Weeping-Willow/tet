package rates

import (
	"context"
	"encoding/xml"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/Weeping-Willow/tet/internal/config"
	"github.com/Weeping-Willow/tet/internal/objects"
	"github.com/pkg/errors"
)

type ecbRSSFetcher struct {
	url    string
	client *http.Client

	currencyFinderRegex *regexp.Regexp
}

func NewEcbRssFetcher(cfg config.Config) Fetcher {
	return &ecbRSSFetcher{
		url: cfg.ExternalServices.EcbRssURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		currencyFinderRegex: regexp.MustCompile(`([A-Z]{3})\s+([\d.]+)`),
	}
}

func (e ecbRSSFetcher) GetCurrencyRate(ctx context.Context, currencyCode string) (objects.CurrencyRate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.url, nil)
	if err != nil {
		return objects.CurrencyRate{}, err
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return objects.CurrencyRate{}, errors.Wrap(err, "fetching currency")
	}

	if resp.StatusCode != http.StatusOK {
		return objects.CurrencyRate{}, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.Body == nil {
		return objects.CurrencyRate{}, errors.New("empty response body")
	}

	defer resp.Body.Close()

	var rss XMLEcbRss
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return objects.CurrencyRate{}, errors.Wrap(err, "decoding xml")
	}

	c := objects.CurrencyRate{
		Currency: currencyCode,
		DayRates: make([]objects.CurrencyRateDay, 0, len(rss.Channel.Items)),
	}

	for _, item := range rss.Channel.Items {
		matches := e.currencyFinderRegex.FindAllStringSubmatch(item.Description, -1)

		for _, match := range matches {
			if len(match) != 3 {
				continue
			}

			if match[1] != currencyCode {
				continue
			}

			rate := match[2]
			rateConverted, err := strconv.ParseFloat(rate, 64)
			if err != nil {
				return objects.CurrencyRate{}, errors.Wrap(err, "parsing rate")
			}

			parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				return objects.CurrencyRate{}, errors.Wrap(err, "parsing date")
			}

			c.DayRates = append(c.DayRates, objects.CurrencyRateDay{
				Rate: rateConverted,
				Date: parsedTime,
			})

			break
		}
	}

	if len(c.DayRates) == 0 {
		return c, errors.Errorf("currency not found: %s", currencyCode)
	}

	return c, nil
}
