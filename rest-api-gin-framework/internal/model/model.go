package model

type BookCreateDTO struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookUpdateDTO struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookResponseDTO struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
