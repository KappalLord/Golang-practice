package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
}

type Rate struct {
	Id    int
	Name  string
	Price float64
}

func getExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	currency := r.URL.Query().Get("currency")

	connStr := "user=postgres password=password dbname=postgres sslmode=disable host=localhost port=5432"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row := db.QueryRow("select * from rates where rate_name=$1", currency)
	DBrate := Rate{}
	err = row.Scan(&DBrate.Id, &DBrate.Name, &DBrate.Price)

	switch err {
	case sql.ErrNoRows:
		resp, err := http.Get("http://api.exchangeratesapi.io/v1/latest?access_key=986cdb6c4c9a47fabf369d77b6b82154")

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var rates ExchangeRates
		if err := json.Unmarshal(body, &rates); err != nil {
			log.Fatal(err)
		}

		rate, exists := rates.Rates[currency]
		if !exists {
			log.Fatal(err)
		}

		result, err := db.Exec("insert into rates (rate_name, rate_price) values ($1, $2)", currency, rate)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)

		response := map[string]float64{
			"rate": rate,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		fmt.Println("Case with API")

	case nil:
		response := map[string]float64{
			"rate": DBrate.Price,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		fmt.Println("Case with DB")
	default:
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", getExchangeRateHandler)

	fmt.Println("Сервер запущен на порту 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
