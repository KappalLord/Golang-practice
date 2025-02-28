package main

import (
	"apiproject/configs"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	memorycache "github.com/maxchagin/go-memorycache-example"
)

func init() {
	// loads values from .env into the system
	err := godotenv.Load("userData.env")
	if err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf := configs.New()
	cache := memorycache.New(5*time.Minute, 10*time.Minute)

	fmt.Println(conf.User)

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s", conf.User, conf.Password, conf.Dbname, conf.Sslmode, conf.Host, conf.Port)

	database, err := OpenDB(connStr)
	if err != nil {
		log.Println(err)
	}
	defer database.db.Close()

	client := NewExchangeRateClient(conf.ApiKey)

	http.HandleFunc("/", getExchangeRateHandler(client, database, cache))
	fmt.Println("Сервер запущен на порту 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
