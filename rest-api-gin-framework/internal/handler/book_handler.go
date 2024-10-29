package handler

import (
	"github/dev-toolkit-go/rest-api-gin-framework/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct{}

func NewBookHandler() *BookHandler {
	return &BookHandler{}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	bookResponse := model.BookResponse{
		ID:     1,
		Title:  book.Title,
		Author: book.Author,
	}

	c.JSON(http.StatusCreated, bookResponse)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid book ID"})
		return
	}

	bookResponse := model.BookResponse{
		ID:     uint(id),
		Title:  "Sample Book",
		Author: "Author Name",
	}

	c.JSON(http.StatusOK, bookResponse)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid book ID"})
		return
	}

	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	bookResponse := model.BookResponse{
		ID:     uint(id),
		Title:  book.Title,
		Author: book.Author,
	}

	c.JSON(http.StatusOK, bookResponse)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	_, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid book ID"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	books := []model.BookResponse{
		{ID: 1, Title: "Sample Book 1", Author: "Author 1"},
		{ID: 2, Title: "Sample Book 2", Author: "Author 2"},
	}

	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetTopRatedBooks(c *gin.Context) {
	topRatedBooks := []model.BookResponse{
		{ID: 1, Title: "Top Rated Book 1", Author: "Author 1"},
		{ID: 2, Title: "Top Rated Book 2", Author: "Author 2"},
	}

	c.JSON(http.StatusOK, topRatedBooks)
}

type ErrorResponse struct {
	Message string `json:"message"`
}
