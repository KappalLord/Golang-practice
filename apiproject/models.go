package main

import (
	"database/sql"
	"time"
)

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
}

type Rate struct {
	Id    int
	Name  string
	Price float64
}

type ExchangeRateApiClient struct {
	apiURL string
}

type Database struct {
	db *sql.DB
}
type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}
