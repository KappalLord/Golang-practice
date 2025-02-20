package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func getExchangeRateHandler(client *ExchangeRateApiClient, database *Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currency := r.URL.Query().Get("currency")

		DBrate, err := database.GetRate(currency)

		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Database error", http.StatusNotAcceptable)
			return
		}

		if err == sql.ErrNoRows {
			rate, err := client.GetExchangeRateApi(currency)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotAcceptable)
				return
			}
			err = database.SaveRate(currency, rate)
			if err != nil {
				log.Println(err)
				return
			}

			response := map[string]float64{
				"rate": rate,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			fmt.Println("Case with API")
		} else {
			response := map[string]float64{
				"rate": DBrate.Price,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			fmt.Println("Case with DB")
		}

	}
}
