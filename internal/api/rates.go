package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Weeping-Willow/tet/internal/utils"
	"github.com/pkg/errors"
)

type ErrorResponse struct {
	ErrorMsg string `json:"error_msg"`
}

type LatestExchangeRatesResponse struct {
	Rates []LatestExchangeRate `json:"rates,omitempty"`
}

type LatestExchangeRate struct {
	CurrencyCode  string  `json:"currency_code"`
	Rate          float64 `json:"rate"`
	EffectiveDate string  `json:"effective_date"`
}

type ExchangeRateHistoryResponse struct {
	CurrencyCode string                `json:"currency_code"`
	Rates        []ExchangeRateHistory `json:"rates,omitempty"`
}

type ExchangeRateHistory struct {
	Rate       float64 `json:"rate"`
	UpdateDate string  `json:"update_date"`
}

func (a *API) latestExchangeRatesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := a.rateService.GetLatestRates(r.Context())
		if err != nil {
			utils.LoggerFromContext(r.Context()).Error(err.Error())
			a.errorResponse(r.Context(), w, http.StatusInternalServerError, err.Error())

			return
		}

		response := LatestExchangeRatesResponse{
			Rates: make([]LatestExchangeRate, 0, len(res)),
		}

		for _, rate := range res {
			if len(rate.DayRates) == 0 {
				continue
			}

			latestDayRate := rate.DayRates[0]
			response.Rates = append(response.Rates, LatestExchangeRate{
				CurrencyCode:  rate.Currency,
				Rate:          latestDayRate.Rate,
				EffectiveDate: latestDayRate.Date.Format(time.DateOnly),
			})
		}

		a.sendResponse(r.Context(), w, http.StatusOK, response)
	})
}

func (a *API) currencyExchangeRateHistoryHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		if code == "" {
			a.errorResponse(r.Context(), w, http.StatusBadRequest, "currency name is required")
		}

		if len(code) != 3 {
			a.errorResponse(r.Context(), w, http.StatusBadRequest, "currency name must be 3 characters long")
		}

		code = strings.ToUpper(code)

		ctx := utils.ContextWithLogging(r.Context(), utils.LoggerFromContext(r.Context()).With("code", code))

		res, err := a.rateService.GetCurrencyHistory(ctx, code)
		if err != nil {
			utils.LoggerFromContext(r.Context()).Error(err.Error())
			a.errorResponse(r.Context(), w, http.StatusInternalServerError, err.Error())

			return
		}

		response := ExchangeRateHistoryResponse{
			CurrencyCode: res.Currency,
			Rates:        make([]ExchangeRateHistory, 0, len(res.DayRates)),
		}

		for _, dayRate := range res.DayRates {
			response.Rates = append(response.Rates, ExchangeRateHistory{
				Rate:       dayRate.Rate,
				UpdateDate: dayRate.Date.Format(time.DateOnly),
			})
		}

		a.sendResponse(r.Context(), w, http.StatusOK, response)
	})
}

func (a *API) errorResponse(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	a.sendResponse(ctx, w, status, ErrorResponse{ErrorMsg: msg})
}

func (a *API) sendResponse(ctx context.Context, w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response, err := json.Marshal(payload)
	if err != nil {
		err = errors.Wrap(err, "marshal response")
		utils.LoggerFromContext(ctx).Error(err.Error())

		_, _ = w.Write([]byte(`{"error_msg":"failed to marshal response"}`))
	}

	_, _ = w.Write(response)
}
