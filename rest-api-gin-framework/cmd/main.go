package main

import (
	"github/dev-toolkit-go/rest-api-gin-framework/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()

	log.Println("Starting server on http://localhost:8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func setupRouter() *gin.Engine {
	bookHandler := handler.NewBookHandler()

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		// Bookstore API Endpoints (CRUD + custom)
		v1.POST("/books", bookHandler.CreateBook)
		v1.GET("/books/:id", bookHandler.GetBook)
		v1.PUT("/books/:id", bookHandler.UpdateBook)
		v1.DELETE("/books/:id", bookHandler.DeleteBook)
		v1.GET("/books", bookHandler.ListBooks)
		v1.GET("/books/top-rated", bookHandler.GetTopRatedBooks)
	}

	return router
}
