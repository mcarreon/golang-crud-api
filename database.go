package main

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func OpenConnection() *sql.DB {
	connStr := os.Getenv("DB_CONN")

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	return db
}
