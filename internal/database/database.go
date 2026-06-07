package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к базе данных: %w", err)
	}

	db.SetMaxOpenConns(20) // установить максимальное количество соединений
	db.SetMaxIdleConns(5)  // установить максимальное количество ПРОСТАИВАЮЩИХ соединений (остаются открытыми и готовы быть переиспользованы)

	return db, nil
}
