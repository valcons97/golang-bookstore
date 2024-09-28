package router

import (
	"bookstore/internal/book/handler"
	"bookstore/internal/book/repository"
	"bookstore/internal/book/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func BookRouter(router *gin.Engine, db *sql.DB) {
	repo := repository.NewBookRepository(db) // Create the repository
	svc := service.NewBookService(repo)      // Create the service
	handler := handler.NewBookHandler(svc)   // Create the handler

	// Define the routes
	router.GET("/book", handler.GetBooks)
	router.GET("/book/:id", handler.GetBookById)
	router.POST("/book/create", handler.CreateBook)
	router.POST("/book/update", handler.UpdateBook)
}
