package objects

import "time"

type CurrencyRate struct {
	Currency string
	DayRates []CurrencyRateDay
}

type CurrencyRateDay struct {
	Rate float64
	Date time.Time
}
