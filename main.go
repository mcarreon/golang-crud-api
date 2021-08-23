package main

import (
	"log"
	"net/http"
	"time"
)

var timePlaceholder = time.Date(
	2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

type InMemoryBookStore struct {
	books []Book
}

func (i *InMemoryBookStore) GetBooks() []Book {
	books := []Book{
		{"Test", "John", timePlaceholder, "Publishers", 5, "CheckedIn"},
		{"Test2", "Jill", timePlaceholder, "Publishers", 3, "CheckedOut"},
	}

	return books
}

func (i *InMemoryBookStore) SaveBook(book Book) {}

func (i *InMemoryBookStore) DeleteBook(title string) {}

func (i *InMemoryBookStore) UpdateBook(title string, fields map[string]interface{}) {}

func (i *InMemoryBookStore) GetBook(title string) Book {
	return Book{}
}

func NewInMemoryBookStore() *InMemoryBookStore {
	return &InMemoryBookStore{}
}

func main() {
	server := NewBookServer(NewInMemoryBookStore())
	log.Fatal(http.ListenAndServe(":3000", server))
}
