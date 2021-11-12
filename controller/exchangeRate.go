package controller

import (
	"encoding/json"
	"exchange-api/config"
	"exchange-api/service"
	"net/http"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/gorilla/mux"
)

type ExchangeRateController interface {
	GetExchangeRate(rw http.ResponseWriter, req *http.Request)
}

type exchangeRateController struct {
	exchangeRateService service.ExchangeRateService
	cache               ttlcache.SimpleCache
}

func NewExchangeRateController(exchangeRateService service.ExchangeRateService, config *config.ExchangeRateControllerConfigurations) ExchangeRateController {
	var cache ttlcache.SimpleCache = ttlcache.NewCache()
	cache.SetTTL(time.Duration(time.Duration(config.CacheDurationMinutes) * time.Minute))
	return &exchangeRateController{
		exchangeRateService: exchangeRateService,
		cache:               cache,
	}
}

func (c *exchangeRateController) GetExchangeRate(rw http.ResponseWriter, req *http.Request) {
	name := mux.Vars(req)["name"]
	if name != "EUR" && name != "USD" && name != "GBP" {
		ResponseCreator(rw, http.StatusBadRequest, nil, "You must provide a valid currency name")
		return
	}
	res, err := c.cache.Get(name)
	if err == nil {
		ResponseCreator(rw, http.StatusOK, res, "")
		return
	}
	exchangeRate, err := c.exchangeRateService.GetBestExchangeRate(name)
	if err != nil {
		ResponseCreator(rw, http.StatusInternalServerError, nil, err.Error())
		return
	}
	ExchangeRateResponse := ExchangeRateResponse{
		Value:    exchangeRate,
		Currency: name,
	}
	c.cache.Set(name, ExchangeRateResponse)
	ResponseCreator(rw, http.StatusOK, ExchangeRateResponse, "")
}

func ResponseCreator(rw http.ResponseWriter, code int, data interface{}, errMessage string) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	response := Response{
		Code: code,
		Err:  errMessage,
		Data: data,
	}
	j, err := json.Marshal(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Write(j)
}
