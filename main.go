package main

import (
	"log"
	"net/http"
)

type InMemoryBookStore struct {
	books []Book
}

func (i *InMemoryBookStore) GetBooks() []Book {
	books := []Book{
		{"Test", "John", "Publishers", 5, "CheckedIn"},
		{"Test2", "Jill", "Publishers", 3, "CheckedOut"},
	}

	return books
}

func (i *InMemoryBookStore) SaveBook(book Book) {}

func NewInMemoryBookStore() *InMemoryBookStore {
	return &InMemoryBookStore{}
}

func main() {
	server := NewBookServer(NewInMemoryBookStore())
	log.Fatal(http.ListenAndServe(":3000", server))
}
