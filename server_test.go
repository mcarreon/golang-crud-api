package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubBookStore struct {
	books []Book
}

func (s *StubBookStore) GetBooks() []Book {
	return s.books
}

func (s *StubBookStore) SaveBook(book Book) {
	s.books = append(s.books, book)
}

func TestGETBooks(t *testing.T) {
	t.Run("it returns 200 on /books", func(t *testing.T) {
		store := StubBookStore{}
		server := NewBookServer(&store)

		request := newGetBooksRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})

	// t.Run("it returns all books as JSON", func(t *testing.T) {
	// 	testBooks := []Book{
	// 		{"Test", "John", "Publishers", 5, "CheckedIn"},
	// 		{"Test2", "Jill", "Publishers", 3, "CheckedOut"},
	// 	}

	// 	store := StubBookStore{testBooks}
	// 	server := NewBookServer(&store)

	// 	request := newGetBooksRequest()
	// 	response := httptest.NewRecorder()

	// 	server.ServeHTTP(response, request)

	// 	got := getBooksFromResponse(t, response.Body)

	// 	assertStatus(t, response.Code, http.StatusOK)
	// 	assertBooks(t, got, testBooks)
	// })
}

func TestPOSTBook(t *testing.T) {
	t.Run("should create new book and get status 200", func(t *testing.T) {
		store := StubBookStore{[]Book{}}
		server := NewBookServer(&store)

		var jsonStr = []byte(`{"title": "Test Book", "author": "John", "publisher": "publisher", "rating": 5, "status": "checkedIn"}`)

		request := newPostBookRequest(jsonStr)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBooksLen(t, len(store.books), 1)
	})

	t.Run("send JSON with empty fields and get status 400", func(t *testing.T) {
		store := StubBookStore{[]Book{}}
		server := NewBookServer(&store)

		var jsonStr = []byte(`{"title": "", "author": "John", "publisher": "publisher", "rating": 5, "status": "checkedIn"}`)

		request := newPostBookRequest(jsonStr)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertBooksLen(t, len(store.books), 0)
	})

	t.Run("send bad JSON and get status 422", func(t *testing.T) {
		store := StubBookStore{[]Book{}}
		server := NewBookServer(&store)

		var jsonStr = []byte(`{"title": "Test Book, "author": "John"}`)

		request := newPostBookRequest(jsonStr)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnprocessableEntity)
		assertBooksLen(t, len(store.books), 0)
	})

}

func TestPUTBook(t *testing.T) {
	t.Run("should recieve 200 status", func(t *testing.T) {
		store := StubBookStore{[]Book{}}
		server := NewBookServer(&store)

		request := newPutBookRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}

func getBooksFromResponse(t testing.TB, body io.Reader) (book []Book) {
	t.Helper()

	book, err := NewBooks(body)

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

func newPostBookRequest(jsonStr []byte) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/book", bytes.NewBuffer(jsonStr))

	return request
}

func newGetBooksRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/books", nil)

	return request
}

func newPutBookRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodPut, "/books/test", nil)

	return request
}
