package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
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

func getBooksFromResponse(t testing.TB, body io.Reader) (book []Book) {
	t.Helper()

	book, err := NewBooks(body)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into books, %v", body, err)
	}

	return book
}

func getBookFromResponse(t testing.TB, body io.Reader) Book {
	t.Helper()

	var book Book

	err := json.NewDecoder(body).Decode(&book)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into books, %v", body, err)
	}

	return book
}

func NewBooks(rdr io.Reader) ([]Book, error) {
	var book []Book
	err := json.NewDecoder(rdr).Decode(&book)

	if err != nil {
		err = fmt.Errorf("problem parsing book, %v", err)
	}

	return book, err
}

// Universal request for GET/DEL specific book
func newBookRequest(method, title string) *http.Request {
	request, _ := http.NewRequest(method, fmt.Sprintf("/books/%s", title), nil)

	return request
}

func newGetBooksRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/books", nil)

	return request
}

func newPostBookRequest(jsonStr []byte) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/book", bytes.NewBuffer(jsonStr))

	return request
}

func newPutBookRequest(title string, jsonStr []byte) *http.Request {
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/books/%s", title), bytes.NewBuffer(jsonStr))

	return request
}
