package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func OpenConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var envData map[string]string
	envData, envError := godotenv.Read()

	if envError != nil {
		log.Fatalf("Error loading.env file")
	}

	psqlCreds := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", envData["HOST"], envData["PORT"], envData["USER"], envData["PASSWORD"], envData["DBNAME"])

	db, err := sql.Open("postgres", psqlCreds)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	return db
}
