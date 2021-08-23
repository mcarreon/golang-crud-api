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
	return (book.Status == "CheckedIn" || book.Status == "CheckedOut")
}

func DecodeBook(r *http.Request) (Book, error) {
	book := Book{}
	err := json.NewDecoder(r.Body).Decode(&book)
	return book, err
}

func DecodeUpdateFields(r *http.Request) (map[string]interface{}, error) {
	var fields map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&fields)

	return fields, err
}

func GetIndexOfStruct(books []Book, title string) int {
	var index int

	for i, item := range books {
		if item.Title == title {
			index = i
		}
	}

	return index
}
