package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	memorycache "github.com/maxchagin/go-memorycache-example"
)

func getExchangeRateHandler(client *ExchangeRateApiClient, database *Database, cache *memorycache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currency := r.URL.Query().Get("currency")

		casheValue, exists := cache.Get(currency)

		if !exists {
			fmt.Println("No exchange rate found in cashe")

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
				cache.Set(currency, rate, 5*time.Minute)

				response := map[string]float64{
					"rate": rate,
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				fmt.Println("Case with API")
			} else {
				cache.Set(currency, DBrate.Price, 5*time.Minute)
				response := map[string]float64{
					"rate": DBrate.Price,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				fmt.Println("Case with DB")
			}
		} else {
			response := map[string]float64{
				"rate": casheValue.(float64),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			fmt.Println("Case with cashe")

		}

	}
}
