package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func NewExchangeRateClient(apiURL string) *ExchangeRateApiClient {
	return &ExchangeRateApiClient{apiURL: apiURL}
}

func (client *ExchangeRateApiClient) GetExchangeRateApi(currency string) (float64, error) {
	resp, err := http.Get(client.apiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var rates ExchangeRates
	if err := json.Unmarshal(body, &rates); err != nil {
		return 0, err
	}

	rate, exists := rates.Rates[currency]
	if !exists {
		return 0, fmt.Errorf("Currency does not exist")
	}

	return rate, nil
}
