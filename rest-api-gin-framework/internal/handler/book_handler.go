// internal/handler/book_handler.go
package handler

import (
	"net/http"
	"strconv"

	"github/dev-toolkit-go/rest-api-gin-framework/internal/model"

	"github.com/gin-gonic/gin"
)

type BookHandler struct{}

func NewBookHandler() *BookHandler {
	return &BookHandler{}
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var createDTO model.BookCreateDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	// Mock response
	createdBook := model.BookResponseDTO{
		ID:     1,
		Title:  createDTO.Title,
		Author: createDTO.Author,
	}

	c.JSON(http.StatusCreated, createdBook)
}

func (h *BookHandler) GetBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid book ID"})
		return
	}

	// Mock response
	book := model.BookResponseDTO{
		ID:     uint(id),
		Title:  "Sample Book",
		Author: "Author Name",
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid book ID"})
		return
	}

	var updateDTO model.BookUpdateDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	// Mock response
	updatedBook := model.BookResponseDTO{
		ID:     uint(id),
		Title:  updateDTO.Title,
		Author: updateDTO.Author,
	}

	c.JSON(http.StatusOK, updatedBook)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	_, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid book ID"})
		return
	}

	// Mock response
	c.JSON(http.StatusNoContent, nil)
}

func (h *BookHandler) ListBooks(c *gin.Context) {
	// Mock response
	books := []model.BookResponseDTO{
		{ID: 1, Title: "Sample Book 1", Author: "Author 1"},
		{ID: 2, Title: "Sample Book 2", Author: "Author 2"},
	}

	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetTopRatedBooks(c *gin.Context) {
	// Mock response
	topRatedBooks := []model.BookResponseDTO{
		{ID: 1, Title: "Top Rated Book 1", Author: "Author 1"},
		{ID: 2, Title: "Top Rated Book 2", Author: "Author 2"},
	}

	c.JSON(http.StatusOK, topRatedBooks)
}

type ErrorResponse struct {
	Message string `json:"message"`
}
