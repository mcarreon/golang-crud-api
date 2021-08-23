package main

import "time"

type Book struct {
	Title          string    `json:"title"`
	Author         string    `json:"author"`
	Published_Date time.Time `json:"publishDate"`
	Publisher      string    `json:"publisher"`
	Rating         int       `json:"rating"`
	Status         string    `json:"status"`
}
