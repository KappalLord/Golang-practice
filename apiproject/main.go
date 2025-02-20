package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

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
