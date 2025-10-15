package storage

import (
	"context"
	"time"

	"github.com/Weeping-Willow/tet/internal/objects"
	"github.com/Weeping-Willow/tet/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type RateRepository struct {
	db *sqlx.DB
}

type Rate struct {
	ID            int64     `db:"id"`
	CurrencyCode  string    `db:"currency_code"`
	Rate          float64   `db:"rate"`
	EffectiveDate time.Time `db:"effective_date"`
}

func NewRateRepository(db *sqlx.DB) *RateRepository {
	return &RateRepository{db: db}
}

func (r *RateRepository) UpdateRates(ctx context.Context, rates []objects.CurrencyRate) error {
	ratesUpd := make([]Rate, 0, len(rates))

	for _, rate := range rates {
		for _, dayRate := range rate.DayRates {
			ratesUpd = append(ratesUpd, Rate{
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
