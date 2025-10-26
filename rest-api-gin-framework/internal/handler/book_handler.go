package handler

import (
	"net/http"
	"strconv"

	"github/dev-toolkit-go/rest-api-gin-framework/internal/model"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
}

func NewBookHandler() *BookHandler {
	return &BookHandler{}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var book model.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{Message: "Invalid input data: " + err.Error()},
		)
		return
	}

	bookResponse := model.BookResponse{
		ID:     99, // Mock ID
		Title:  book.Title,
		Author: book.Author,
		Rating: book.Rating,
	}

	c.JSON(http.StatusCreated, bookResponse)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{Message: "Invalid book ID format. Must be an integer."},
		)
		return
	}

	bookResponse := model.BookResponse{
		ID:     uint(id),
		Title:  "Mock Book Title",
		Author: "Mock Author Name",
		Rating: 5,
	}

	c.JSON(http.StatusOK, bookResponse)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{Message: "Invalid book ID format. Must be an integer."},
		)
		return
	}

	var bookUpdate model.Book
	if err := c.ShouldBindJSON(&bookUpdate); err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{Message: "Invalid input data: " + err.Error()},
		)
		return
	}

	updatedBook := model.BookResponse{
		ID:     uint(id),
		Title:  bookUpdate.Title,
		Author: bookUpdate.Author,
		Rating: bookUpdate.Rating,
	}

	c.JSON(http.StatusOK, updatedBook)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	// 1. Parse ID
	idStr := c.Param("id")
	_, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{Message: "Invalid book ID format. Must be an integer."},
		)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	books := []model.BookResponse{
		{ID: 101, Title: "Favourite Book 1", Author: "Favourite Author 1 ", Rating: 5},
		{ID: 102, Title: "Favourite Book 1", Author: "Favourite Author 2", Rating: 4},
	}

	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetTopRatedBooks(c *gin.Context) {
	topRatedBooks := []model.BookResponse{
		{ID: 201, Title: "Favourite Book 1", Author: "Favourite Author 1", Rating: 5},
		{ID: 202, Title: "Favourite Book 2", Author: "Favourite Author 2", Rating: 4},
	}

	c.JSON(http.StatusOK, topRatedBooks)
}
