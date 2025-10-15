package storage

import (
	"context"
	"time"

	"github.com/Weeping-Willow/tet/internal/objects"
	"github.com/Weeping-Willow/tet/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Rates struct {
	db *sqlx.DB
}

type RateDB struct {
	ID            int64     `db:"id"`
	CurrencyCode  string    `db:"currency_code"`
	Rate          float64   `db:"rate"`
	EffectiveDate time.Time `db:"effective_date"`
}

func NewRates(db *sqlx.DB) *Rates {
	return &Rates{db: db}
}

func (r *Rates) UpdateRates(ctx context.Context, rates []objects.CurrencyRate) error {
	ratesUpd := make([]RateDB, 0, len(rates))

	for _, rate := range rates {
		for _, dayRate := range rate.DayRates {
			ratesUpd = append(ratesUpd, RateDB{
				CurrencyCode:  rate.Currency,
				Rate:          dayRate.Rate,
				EffectiveDate: dayRate.Date,
			})
		}
	}

	query := `
		INSERT INTO exchange_rates (currency_code, rate, effective_date)
		VALUES (:currency_code, :rate, :effective_date)
		ON DUPLICATE KEY UPDATE rate = VALUES(rate)
	`

	res, err := r.db.NamedExecContext(ctx, query, ratesUpd)
	if err != nil {
		return errors.Wrap(err, "updating rates in database")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "getting rows affected")
	}

	utils.LoggerFromContext(ctx).Info("Rates updated", "rows", rows)

	return nil

}

func (r *Rates) GetLatestRates(ctx context.Context) ([]objects.CurrencyRate, error) {
	query := `
		SELECT currency_code, rate, effective_date
		FROM exchange_rates er
		WHERE effective_date = (
			SELECT MAX(effective_date)
			FROM exchange_rates
			WHERE currency_code = er.currency_code
		)
	`

	var rates []RateDB
	err := r.db.SelectContext(ctx, &rates, query)
	if err != nil {
		return nil, errors.Wrap(err, "selecting latest rates from db")
	}

	result := make([]objects.CurrencyRate, 0, len(rates))

	for _, rate := range rates {
		result = append(result, objects.CurrencyRate{
			Currency: rate.CurrencyCode,
			DayRates: []objects.CurrencyRateDay{
				{
					Rate: rate.Rate,
					Date: rate.EffectiveDate,
				},
			},
		})
	}

	return result, nil
}

func (r *Rates) GetCurrencyHistory(ctx context.Context, currencyCode string) (objects.CurrencyRate, error) {
	query := `
		SELECT currency_code, rate, effective_date
		FROM exchange_rates
		WHERE currency_code = ?
		ORDER BY effective_date DESC
	`

	var rates []RateDB
	err := r.db.SelectContext(ctx, &rates, query, currencyCode)
	if err != nil {
		return objects.CurrencyRate{}, errors.Wrapf(err, "selecting rates for currency %s from db", currencyCode)
	}

	if len(rates) == 0 {
		return objects.CurrencyRate{}, errors.Errorf("no rates found for currency %s", currencyCode)
	}

	result := objects.CurrencyRate{
		Currency: currencyCode,
		DayRates: make([]objects.CurrencyRateDay, 0, len(rates)),
	}

	for _, rate := range rates {
		result.DayRates = append(result.DayRates, objects.CurrencyRateDay{
			Rate: rate.Rate,
			Date: rate.EffectiveDate,
		})
	}

	return result, nil
}
