package api

import "net/http"

func (a *API) latestExchangeRatesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`NOT IMPLEMENTED`))
	})
}

func (a *API) currencyExchangeRateHistoryHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`NOT IMPLEMENTED`))
	})
}
