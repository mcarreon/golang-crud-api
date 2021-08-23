package main

import (
	"encoding/json"
	"net/http"
)

const jsonContentType = "application/json"

type BookStore interface {
	GetBooks() []Book
	SaveBook(book Book)
}

type BookServer struct {
	store BookStore
	http.Handler
}

func NewBookServer(store BookStore) *BookServer {
	b := &BookServer{
		store,
		http.NewServeMux(),
	}

	b.store = store

	router := http.NewServeMux()
	router.Handle("/books", http.HandlerFunc(b.booksHandler))
	router.Handle("/book", http.HandlerFunc(b.bookHandler))

	b.Handler = router

	return b
}

// Handles GET/PUT/DELETE requests for specific books
func (b *BookServer) bookHandler(w http.ResponseWriter, r *http.Request) {
	book := Book{}

	err := json.NewDecoder(r.Body).Decode(&book)

	// If unable to parse bad JSON, 422
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	hasGoodValues := CheckNotEmpty(book)
	hasGoodStatus := ValidStatus(book)

	// If empty values or bad status, 400
	if !hasGoodValues || !hasGoodStatus {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b.store.SaveBook(book)
}

// Handles GET request for all books
func (b *BookServer) booksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)

	testBooks := []Book{
		{"Test", "John", "Publishers", 5, "CheckedIn"},
		{"Test2", "Jill", "Publishers", 3, "CheckedOut"},
	}

	json.NewEncoder(w).Encode(testBooks)
	//json.NewEncoder(w).Encode(b.store.GetBooks())
}

func (b *BookServer) postBook(w http.ResponseWriter, r *http.Request) {

}
