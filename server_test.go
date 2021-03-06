package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type StubBookStore struct {
	books []Book
}

var timePlaceholder = time.Date(
	2009, 11, 17, 20, 34, 58, 651387237, time.UTC)

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
	index := GetIndexOfStruct(s.books, title)

	s.books[index] = s.books[len(s.books)-1]
	s.books = s.books[:len(s.books)-1]
}

func (s *StubBookStore) UpdateBook(title string, fields map[string]interface{}) {
	index := GetIndexOfStruct(s.books, title)

	for key, field := range fields {
		switch key {
		case "title":
			s.books[index].Title = field.(string)
		case "author":
			s.books[index].Author = field.(string)
		case "publisher":
			s.books[index].Publisher = field.(string)
		case "publishedDate":
			parsedTime, _ := time.Parse("2006-01-02T15:04:05.000Z", field.(string))
			s.books[index].Published_Date = parsedTime
		case "rating":
			floatNum := field.(float64)
			s.books[index].Rating = int(floatNum)
		case "status":
			s.books[index].Status = field.(string)
		}
	}
}

func TestGETBooks(t *testing.T) {
	t.Run("it returns all books as JSON and returns 200", func(t *testing.T) {
		testBooks := []Book{
			{"Test", "John", timePlaceholder, "Publishers", 5, "CheckedIn"},
			{"Test2", "Jill", timePlaceholder, "Publishers", 3, "CheckedOut"},
		}

		store := StubBookStore{testBooks}
		server := NewBookServer(&store)

		request := newGetBooksRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getBooksFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertBooks(t, got, testBooks)
	})

	t.Run("should not grab any books, and send back 404", func(t *testing.T) {
		store := StubBookStore{}
		server := NewBookServer(&store)

		request := newGetBooksRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestGETBook(t *testing.T) {
	testBook := []Book{
		{"Test", "John", timePlaceholder, "Publishers", 5, "CheckedIn"},
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

		var jsonStr = []byte(`{"title": "Test Book", "author": "John", "publishDate": "2005-06-13T04:40:51Z", "publisher": "publisher", "rating": 3, "status": "CheckedIn"}`)

		request := newPostBookRequest(jsonStr)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertBooksLen(t, len(store.books), 1)
	})

	t.Run("send JSON with empty fields and get status 400", func(t *testing.T) {
		store := StubBookStore{[]Book{}}
		server := NewBookServer(&store)

		var jsonStr = []byte(`{"title": "", "author": "John", "publisher": "publisher", "rating": 5, "status": "CheckedIn"}`)

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
			{"Test", "John", timePlaceholder, "Publishers", 3, "CheckedIn"},
		}
		targetTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2006-01-02T15:04:05.000Z")
		targetBook := Book{"Test", "John", targetTime, "Publishers", 3, "CheckedOut"}

		store := StubBookStore{testBook}
		server := NewBookServer(&store)

		var jsonStr = []byte(`{"rating": 3, "publishedDate": "2006-01-02T15:04:05.000Z", "status": "CheckedOut"}`)

		request := newPutBookRequest("Test", jsonStr)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := store.GetBook("Test")

		assertStatus(t, response.Code, http.StatusOK)
		assertBook(t, got, targetBook)
	})

	t.Run("should fail to find book and recieve 404 status", func(t *testing.T) {
		store := StubBookStore{}
		server := NewBookServer(&store)

		var jsonStr = []byte(`{"rating": 3, "status": "CheckedOut"}`)

		request := newPutBookRequest("Test", jsonStr)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("should fail to parse input and recieve 422 status", func(t *testing.T) {
		store := StubBookStore{}
		server := NewBookServer(&store)

		var jsonStr = []byte(`{"rating": 3, "status": "CheckedOut}`)

		request := newPutBookRequest("Test", jsonStr)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnprocessableEntity)
	})

	t.Run("should fail due to bad rating and recieve 422 status", func(t *testing.T) {
		testBook := []Book{
			{"Testing", "John", timePlaceholder, "Publishers", 3, "CheckedIn"},
		}
		store := StubBookStore{testBook}
		server := NewBookServer(&store)

		var jsonStr = []byte(`{"rating": 4, "status": "CheckedOut"}`)

		request := newPutBookRequest("Testing", jsonStr)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnprocessableEntity)
	})
}

func TestDELBook(t *testing.T) {
	testBook := []Book{
		{"Test", "John", timePlaceholder, "Publishers", 3, "CheckedIn"},
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
