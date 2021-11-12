package controller

type Response struct {
	Code int         `json:"code"`
	Err  string      `json:"err"`
	Data interface{} `json:"data"`
}

type ExchangeRateResponse struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}
