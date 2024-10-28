package model

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookResponse struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
