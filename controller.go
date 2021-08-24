package main

import (
	"fmt"
	"strings"
)

const releventCols = "title, author, published_date, publisher, rating, status"

type PostgresStore struct{}

func NewPostgresStore() *PostgresStore {
	return &PostgresStore{}
}

func (p *PostgresStore) GetBooks() []Book {
	db := OpenConnection()
	defer db.Close()

	fmt.Printf("Getting all books \n")

	var books []Book

	rows, err := db.Query(fmt.Sprintf("select %v from books", releventCols))

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

	fmt.Printf("Getting book with title = %v \n", title)

	rows, err := db.Query(fmt.Sprintf("select %v from books where title = '%v'", releventCols, title))

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

	fmt.Printf("Deleting book with title = %v \n", title)

	_, err := db.Query(fmt.Sprintf("delete from books where title = '%v'", title))

	if err != nil {
		panic(err)
	}
}

//TODO: double check date, try to clean
func (p *PostgresStore) SaveBook(book Book) {
	db := OpenConnection()
	defer db.Close()

	fmt.Printf("Posting book with title = %v \n", book.Title)

	_, err := db.Query(fmt.Sprintf("insert into books(%v) values('%v', '%v', '%v', '%v', '%v', '%v')", releventCols, book.Title, book.Author, "2005-06-13T04:40:51Z", book.Publisher, book.Rating, book.Status))

	if err != nil {
		panic(err)
	}
}

func (p *PostgresStore) UpdateBook(title string, fields map[string]interface{}) {
	db := OpenConnection()
	defer db.Close()

	fmt.Printf("Updating book with title = %v \n", title)

	// Temporary slices for building query string and args
	var updateString []string
	var updateArgs []interface{}
	counter := 1

	for key, value := range fields {
		updateString = append(updateString, fmt.Sprintf("%v=$%v", key, counter))
		updateArgs = append(updateArgs, value)
		counter++
	}

	queryStatement := fmt.Sprintf("update books set %s where title = '%s'", strings.Join(updateString, ","), title)

	_, err := db.Exec(queryStatement, updateArgs...)

	if err != nil {
		panic(err)
	}
}
