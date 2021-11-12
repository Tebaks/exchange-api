package client

import (
	"encoding/json"
	"errors"
	"exchange-api/model"
	"fmt"
	"io"
	"net/http"
)

const (
	maxReadBytes      = 10 * 1024 * 1024
	maxReadBytesOnErr = 4 * 1024
)

type ExchangeRateClient interface {
	GetExchangeRate(url string, name string) (float64, error)
}
type exchangeRateClient struct {
	hc *http.Client
}

func NewExchangeRateClient() ExchangeRateClient {
	return &exchangeRateClient{hc: http.DefaultClient}
}

func (c *exchangeRateClient) do(dst interface{}, req *http.Request, expectedStatusCode int) error {
	resp, err := c.hc.Do(req)
	if err != nil {
		return errors.New("failed to do request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxReadBytesOnErr))
		if err != nil {
			return fmt.Errorf("unable to read HTTP response for status code: %s (expected: %d) err: %w", resp.Status, expectedStatusCode, err)
		}

		return fmt.Errorf("remote service returns unexpected response: %s - %s", resp.Status, string(respBody))
	}

	if dst != nil {
		dec := json.NewDecoder(io.LimitReader(resp.Body, maxReadBytes))

		if err := dec.Decode(dst); err != nil {
			return fmt.Errorf("unable to parse JSON response of the remote service err: %w", err)
		}
	}

	return nil
}

func (c *exchangeRateClient) GetExchangeRate(name string, url string) (float64, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	var response ExchangeRatesResponse

	if err := c.do(&response, req, http.StatusOK); err != nil {
		return 0, err
	}

	exchangeRates := model.ExchangeRates{
		EUR: response.EUR,
		USD: response.USD,
		GBP: response.GBP,
	}

	if name == "EUR" {
		return exchangeRates.EUR, nil
	}

	if name == "USD" {
		return exchangeRates.USD, nil
	}

	if name == "GBP" {
		return exchangeRates.GBP, nil
	}

	return 0, errors.New("currency not found")
}
