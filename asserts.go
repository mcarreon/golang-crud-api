package main

import (
	"reflect"
	"testing"
)

// asserts status code
func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is incorrect, got %q, want %q", got, want)
	}
}

func assertBook(t testing.TB, got, want Book) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

// asserts against # of books in store
func assertBooksLen(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("got %d books, expecting %d books", got, want)
	}
}
