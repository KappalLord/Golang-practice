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

type ExchangeRateApiClient struct {
	apiURL string
}

type Database struct {
	db *sql.DB
}

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
		return 0, fmt.Errorf("Валюты не найдено")
	}

	return rate, nil
}

func OpenDB(connStr string) (*Database, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (database *Database) GetRate(currency string) (Rate, error) {
	row := database.db.QueryRow("SELECT * FROM rates WHERE rate_name=$1", currency)
	DBrate := Rate{}
	err := row.Scan(&DBrate.Id, &DBrate.Name, &DBrate.Price)
	if err != nil {
		return DBrate, err
	}
	return DBrate, nil
}

func (database *Database) SaveRate(currency string, rate float64) error {
	_, err := database.db.Exec("INSERT INTO rates (rate_name, rate_price) VALUES ($1, $2)", currency, rate)
	return err
}

func getExchangeRateHandler(client *ExchangeRateApiClient, database *Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currency := r.URL.Query().Get("currency")

		DBrate, err := database.GetRate(currency)

		if err != nil && err != sql.ErrNoRows {
			log.Println(err)
		}
		
		switch DBrate.Price {
		case 0:
			rate, err := client.GetExchangeRateApi(currency)
			if err != nil {
				log.Println(err)
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
			
		default:
			response := map[string]float64{
				"rate": DBrate.Price,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			fmt.Println("Case with DB")

		}
	}
}

func main() {
	connStr := "user=postgres password=password dbname=postgres sslmode=disable host=localhost port=5432"
	database, err := OpenDB(connStr)
	if err != nil {
		log.Println(err)
	}
	defer database.db.Close()

	client := NewExchangeRateClient("http://api.exchangeratesapi.io/v1/latest?access_key=986cdb6c4c9a47fabf369d77b6b82154")

	http.HandleFunc("/", getExchangeRateHandler(client, database))
	fmt.Println("Сервер запущен на порту 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
