package main

import (
	"fmt"
)

const releventCols = "title, author, published_date, publisher, rating, status"

type PostgresStore struct{}

func NewPostgresStore() *PostgresStore {
	return &PostgresStore{}
}

func (p *PostgresStore) GetBooks() []Book {
	db := OpenConnection()
	defer db.Close()

	var books []Book

	rows, err := db.Query(fmt.Sprintf("select %v from go.books", releventCols))

	defer rows.Close()

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var book Book
		rows.Scan(&book.Title, &book.Author, &book.Published_Date, &book.Publisher, &book.Rating, &book.Status)

		books = append(books, book)
	}

	return books
}

func (p *PostgresStore) GetBook(title string) Book {
	db := OpenConnection()
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select %v from go.books where title = '%v'", releventCols, title))

	defer rows.Close()

	if err != nil {
		panic(err)
	}

	var book Book

	for rows.Next() {
		rows.Scan(&book.Title, &book.Author, &book.Published_Date, &book.Publisher, &book.Rating, &book.Status)
	}

	return book
}

func (p *PostgresStore) DeleteBook(title string) {
	db := OpenConnection()
	defer db.Close()

	_, err := db.Query(fmt.Sprintf("delete from go.books where title = '%v'", title))

	if err != nil {
		panic(err)
	}
}

//TODO: double check date, try to clean
func (p *PostgresStore) SaveBook(book Book) {
	db := OpenConnection()
	defer db.Close()

	_, err := db.Query(fmt.Sprintf("insert into go.books(%v) values('%v', '%v', '%v', '%v', '%v', '%v')", releventCols, book.Title, book.Author, "2005-06-13T04:40:51Z", book.Publisher, book.Rating, book.Status))

	if err != nil {
		panic(err)
	}
}

func (p *PostgresStore) UpdateBook(title string, fields map[string]interface{}) {

}
