package main

import (
	"log"
	"net/http"
	"time"
)

var timePlaceholder = time.Date(
	2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

func main() {
	server := NewBookServer(NewPostgresStore())
	log.Fatal(http.ListenAndServe(":3000", server))
}
