package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewBookServer(NewPostgresStore())
	log.Fatal(http.ListenAndServe(":8080", server))
}
