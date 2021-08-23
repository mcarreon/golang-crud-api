package main

type Book struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	Rating    int    `json:"rating"`
	Status    string `json:"status"`
}
