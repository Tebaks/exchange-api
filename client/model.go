package client

type ExchangeRatesResponse struct {
	EUR float64 `json:"EUR"`
	USD float64 `json:"USD"`
	GBP float64 `json:"GBP"`
}
