package model

type Book struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Rating int    `json:"rating" binding:"gte=1,lte=5"`
}

type BookResponse struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
