package service

import (
	"exchange-api/client"
	"exchange-api/config"
)

type ExchangeRateService interface {
	GetBestExchangeRate(name string) (float64, error)
}

type exchangeRateService struct {
	Urls               []string
	exchangeRateClient client.ExchangeRateClient
}

func NewExchangeRateService(exchangeRateClient client.ExchangeRateClient, config *config.ExchangeRateServiceConfigurations) ExchangeRateService {
	return &exchangeRateService{exchangeRateClient: exchangeRateClient,
		Urls: config.Urls}
}

func (s *exchangeRateService) GetBestExchangeRate(name string) (float64, error) {
	resChan := make(chan float64)
	errChan := make(chan error)

	for _, url := range s.Urls {
		go s.GetExchangeRate(name, url, resChan, errChan)
	}

	var bestExcRate float64
	var bestExcRateErr error
	for i := 0; i < len(s.Urls); i++ {
		select {
		case res := <-resChan:
			if bestExcRate == 0 || res < bestExcRate {
				bestExcRate = res
			}
		case err := <-errChan:
			bestExcRateErr = err

		}
		if bestExcRateErr != nil {
			return 0, bestExcRateErr
		}
	}

	return bestExcRate, nil
}

func (s *exchangeRateService) GetExchangeRate(name string, url string, resChan chan<- float64, errChan chan<- error) {
	res, err := s.exchangeRateClient.GetExchangeRate(name, url)
	if err != nil {
		errChan <- err
		return
	}

	resChan <- res
}
