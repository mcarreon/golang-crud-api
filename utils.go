package main

import (
	"encoding/json"
	"net/http"
)

// Checks for empty values
// Not sure how to write this programatically
func CheckNotEmpty(book Book) bool {
	return (book.Title != "" && book.Author != "" && book.Publisher != "" && book.Rating != 0 && book.Status != "")
}

// Checks for correct status
func ValidStatus(book Book) bool {
	return (book.Status == "checkedIn" || book.Status == "checkedOut")
}

func DecodeBook(r *http.Request) (Book, error) {
	book := Book{}
	err := json.NewDecoder(r.Body).Decode(&book)
	return book, err
}
