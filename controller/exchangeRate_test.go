package controller

import (
	"encoding/json"
	"exchange-api/config"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type TestExchangeRateService struct {
	datas map[string]float64
}

func (s *TestExchangeRateService) GetBestExchangeRate(currency string) (float64, error) {
	return s.datas[currency], nil
}

func TestGetExchangeRate(t *testing.T) {
	t.Run("Get Exchange rate with correct name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/currency/EUR", nil)
		vars := map[string]string{
			"name": "EUR",
		}
		req = mux.SetURLVars(req, vars)
		testExchangeRateService := &TestExchangeRateService{
			datas: map[string]float64{
				"EUR": 1.0,
			},
		}
		config := config.ExchangeRateControllerConfigurations{
			CacheDurationMinutes: 0,
		}
		exchangeRateController := NewExchangeRateController(testExchangeRateService, &config)

		w := httptest.NewRecorder()
		exchangeRateController.GetExchangeRate(w, req)

		resp := w.Result()

		body, _ := io.ReadAll(resp.Body)
		var response Response
		err := json.Unmarshal(body, &response)
		if err != nil {
			t.Errorf("Error unmarshalling response: %s", err)
		}
		md, _ := response.Data.(map[string]interface{})
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 1.0, md["value"])
	})
	t.Run("Get Exchange rate with incorrect name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/currency/EUR", nil)
		vars := map[string]string{
			"name": "DOLAR",
		}
		req = mux.SetURLVars(req, vars)
		testExchangeRateService := &TestExchangeRateService{
			datas: map[string]float64{
				"EUR": 1.0,
			},
		}
		config := config.ExchangeRateControllerConfigurations{
			CacheDurationMinutes: 0,
		}
		exchangeRateController := NewExchangeRateController(testExchangeRateService, &config)

		w := httptest.NewRecorder()
		exchangeRateController.GetExchangeRate(w, req)

		resp := w.Result()

		body, _ := io.ReadAll(resp.Body)
		var response Response
		err := json.Unmarshal(body, &response)
		if err != nil {
			t.Errorf("Error unmarshalling response: %s", err)
		}
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, "You must provide a valid currency name", response.Err)
	})
}
