package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const firstBook = `{"title": "Integration Book", "author": "John", "publishDate": "2005-06-13T04:40:51Z", "publisher": "publisher", "rating": 3, "status": "CheckedIn"}`
const updateBook = `{"rating": 1, "status": "CheckedOut"}`

// Теsts creating a book, getting it, updating it, then deleting it
// Also creates a book to ensure nonempty datastore and attempts to grab all books
func TestFullFunctionality(t *testing.T) {
	store := NewPostgresStore()
	server := NewBookServer(store)

	server.ServeHTTP(httptest.NewRecorder(), newPostBookRequest([]byte(firstBook)))

	t.Run("get Integration Book, update it, then delete it", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newBookRequest(http.MethodGet, "Integration Book"))

		// Retrieve created book
		got := getBookFromResponse(t, response.Body)
		assertField(t, got.Title, "Integration Book")
		assertStatus(t, response.Code, http.StatusOK)

		// Update created book
		putResponse := httptest.NewRecorder()
		putGetResponse := httptest.NewRecorder()
		server.ServeHTTP(putResponse, newPutBookRequest("Integration Book", []byte(updateBook)))
		assertStatus(t, putResponse.Code, http.StatusOK)
		server.ServeHTTP(putGetResponse, newBookRequest(http.MethodGet, "Integration Book"))
		got = getBookFromResponse(t, putGetResponse.Body)
		assertField(t, got.Status, "CheckedOut")
		assertStatus(t, putGetResponse.Code, http.StatusOK)

		// Delete created book
		delResponse := httptest.NewRecorder()
		getResponse := httptest.NewRecorder()
		server.ServeHTTP(delResponse, newBookRequest(http.MethodDelete, "Integration Book"))
		assertStatus(t, delResponse.Code, http.StatusOK)
		server.ServeHTTP(getResponse, newBookRequest(http.MethodGet, "Integration Book"))
		assertStatus(t, getResponse.Code, http.StatusNotFound)
	})

	server.ServeHTTP(httptest.NewRecorder(), newPostBookRequest([]byte(firstBook)))

	t.Run("recieve a non-zero amount of books back", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetBooksRequest())

		got := getBooksFromResponse(t, response.Body)

		assertGreaterThanLen(t, len(got), 0)
	})
}
