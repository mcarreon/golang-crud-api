package main

// Checks for empty values
// Not sure how to write this programatically
func CheckNotEmpty(book Book) bool {
	return (book.Title != "" && book.Author != "" && book.Publisher != "" && book.Rating != 0 && book.Status != "")
}

// Checks for correct status
func ValidStatus(book Book) bool {
	return (book.Status == "checkedIn" || book.Status == "checkedOut")
}
