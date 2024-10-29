
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github/dev-toolkit-go/rest-api-gin-framework/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()
	srv := &http.Server{Addr: ":8080", Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	waitForShutdown(srv)
	log.Println("Server exiting")
}

func setupRouter() *gin.Engine {
	bookHandler := handler.NewBookHandler()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/books", bookHandler.CreateBook)
		v1.GET("/books/:id", bookHandler.GetBook)
		v1.PUT("/books/:id", bookHandler.UpdateBook)
		v1.DELETE("/books/:id", bookHandler.DeleteBook)
		v1.GET("/books", bookHandler.ListBooks)
		v1.GET("/books/top-rated", bookHandler.GetTopRatedBooks)
	}
	return router
}

func waitForShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
