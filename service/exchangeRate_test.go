package service

import (
	"errors"
	"exchange-api/config"
	"exchange-api/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testExchangeRateClient struct {
	data map[string]model.ExchangeRates
}

func (c *testExchangeRateClient) GetExchangeRate(name string, url string) (float64, error) {
	res := c.data[url]
	if name == "EUR" {
		return res.EUR, nil
	}
	if name == "USD" {
		return res.USD, nil
	}
	if name == "GBP" {
		return res.GBP, nil
	}

	return 0, errors.New("currency not found")
}

func TestGetBestExchangeRate(t *testing.T) {
	client := &testExchangeRateClient{
		data: map[string]model.ExchangeRates{
			"https://run.mocky.io/v3/f05a9177-5719-47b3-94d9-60ab75e538a7": model.ExchangeRates{
				EUR: 1,
				USD: 1.1,
				GBP: 1.2,
			},
			"https://run.mocky.io/v3/7dcc2ac7-85f0-4032-a421-f8d947e20824": model.ExchangeRates{
				EUR: 0.9,
				USD: 1,
				GBP: 1.1,
			},
			"https://run.mocky.io/v3/377b087e-56bf-4a3f-a2f8-d662b4782705": model.ExchangeRates{
				EUR: 0.8,
				USD: 0.9,
				GBP: 1,
			},
		},
	}
	exchangeRateServiceConfig := config.ExchangeRateServiceConfigurations{
		Urls: []string{"https://run.mocky.io/v3/f05a9177-5719-47b3-94d9-60ab75e538a7", "https://run.mocky.io/v3/7dcc2ac7-85f0-4032-a421-f8d947e20824", "https://run.mocky.io/v3/377b087e-56bf-4a3f-a2f8-d662b4782705"},
	}

	exchangeRateService := NewExchangeRateService(client, &exchangeRateServiceConfig)

	result, err := exchangeRateService.GetBestExchangeRate("EUR")
	assert.NoError(t, err)
	assert.Equal(t, 0.8, result)

	result, err = exchangeRateService.GetBestExchangeRate("USD")
	assert.NoError(t, err)
	assert.Equal(t, 0.9, result)

	result, err = exchangeRateService.GetBestExchangeRate("GBP")
	assert.NoError(t, err)
	assert.Equal(t, 1.0, result)

	result, err = exchangeRateService.GetBestExchangeRate("RUB")
	assert.Error(t, err)
	assert.Equal(t, 0.0, result)
}
