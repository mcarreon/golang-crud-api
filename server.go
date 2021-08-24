package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

const jsonContentType = "application/json"

var timePlaceholder = time.Date(
	2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

type BookStore interface {
	GetBooks() []Book
	GetBook(title string) Book
	SaveBook(book Book)
	UpdateBook(title string, fields map[string]interface{})
	DeleteBook(title string)
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
	router.Handle("/books/", http.HandlerFunc(b.bookHandler))

	b.Handler = router

	return b
}

// Handles GET/PUT/DELETE requests for specific books
func (b *BookServer) bookHandler(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimPrefix(r.URL.Path, "/books/")

	switch r.Method {
	case http.MethodGet:
		b.getBook(w, r, title)
	case http.MethodPut:
		b.updateBook(w, r, title)
	case http.MethodDelete:
		b.deleteBook(w, r, title)
	case http.MethodPost:
		b.postBook(w, r)
	}
}

// Handles GET request for all books and Post request
func (b *BookServer) booksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b.getBooks(w, r)
	case http.MethodPost:
		b.postBook(w, r)
	}
}

// Get all books
func (b *BookServer) getBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	w.Header().Set("content-type", jsonContentType)
	books := b.store.GetBooks()

	if len(books) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(b.store.GetBooks())
}

// POST request functionality for books
func (b *BookServer) postBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	book, err := DecodeBook(r)

	// If unable to parse bad JSON, 422
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	hasGoodValues := CheckNotEmpty(book)
	hasGoodStatus := ValidStatus(book)
	hasGoodRating := ValidRating(book)

	// If empty values or bad status, 400
	if !hasGoodValues || !hasGoodStatus || !hasGoodRating {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b.store.SaveBook(book)
}

// GET request functionality for single book
func (b *BookServer) getBook(w http.ResponseWriter, r *http.Request, title string) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	w.Header().Set("content-type", jsonContentType)
	book := b.store.GetBook(title)

	if (Book{}) == book {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// PUT request functionality
func (b *BookServer) updateBook(w http.ResponseWriter, r *http.Request, title string) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	fields, err := DecodeUpdateFields(r)

	if fields["rating"] != nil && (fields["rating"].(float64) < 0 || fields["rating"].(float64) > 3) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// If unable to process update fields, 422
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	book := b.store.GetBook(title)

	// If unable to find book, 404
	if (Book{}) == book {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b.store.UpdateBook(title, fields)
}

// DEL request functionality
func (b *BookServer) deleteBook(w http.ResponseWriter, r *http.Request, title string) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	book := b.store.GetBook(title)

	if (Book{}) == book {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		b.store.DeleteBook(title)
	}
}
