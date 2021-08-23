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

func (s *StubBookStore) GetBook(title string) Book {
	for _, item := range s.books {
		if item.Title == title {
			return item
		}
	}

	return Book{}
}

func (s *StubBookStore) SaveBook(book Book) {
	s.books = append(s.books, book)
}

func (s *StubBookStore) DeleteBook(title string) {
	var index int

	for i, item := range s.books {
		if item.Title == title {
			index = i
		}
	}

	s.books[index] = s.books[len(s.books)-1]
	s.books = s.books[:len(s.books)-1]
}

func (s *StubBookStore) UpdateBook(title string) {}

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

func TestGETBook(t *testing.T) {
	testBook := []Book{
		{"Test", "John", "Publishers", 5, "CheckedIn"},
	}
	store := StubBookStore{testBook}
	server := NewBookServer(&store)

	t.Run("should get a book and recieve 200 status", func(t *testing.T) {
		request := newBookRequest(http.MethodGet, "Test")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getBookFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertBook(t, got, testBook[0])
	})

	t.Run("should fail to find a book and recieve 404 status", func(t *testing.T) {
		request := newBookRequest(http.MethodGet, "")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
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
	t.Run("should update book and recieve 200 status", func(t *testing.T) {
		testBook := []Book{
			{"Test", "John", "Publishers", 5, "CheckedIn"},
		}
		store := StubBookStore{testBook}
		server := NewBookServer(&store)

		var jsonStr = []byte(`{"rating": 3, "status": "checkedOut"}`)

		request := newPutBookRequest("Test", jsonStr)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}

func TestDELBook(t *testing.T) {
	testBook := []Book{
		{"Test", "John", "Publishers", 5, "CheckedIn"},
	}
	store := StubBookStore{testBook}
	server := NewBookServer(&store)

	t.Run("should delete item and recieve 200 status", func(t *testing.T) {
		request := newBookRequest(http.MethodDelete, "Test")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBooksLen(t, len(store.books), 0)
	})

	t.Run("should fail to find item and recieve 404 status", func(t *testing.T) {
		request := newBookRequest(http.MethodDelete, "Test")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
		assertBooksLen(t, len(store.books), 0)
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
