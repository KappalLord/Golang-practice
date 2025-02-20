package main

import (
	"database/sql"
)

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
