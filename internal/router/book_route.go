package router

import (
	"bookstore/internal/handler"
	"bookstore/internal/repository"
	"bookstore/internal/service"

	"database/sql"

	"github.com/gin-gonic/gin"
)

func BookRouter(router *gin.Engine, db *sql.DB) {
	repo := repository.NewBookRepository(db)
	svc := service.NewBookService(repo)
	handler := handler.NewBookHandler(svc)

	// Define the routes
	router.GET("/book", handler.GetBooks)
	router.GET("/book/:id", handler.GetBookById)
	router.POST("/book/create", handler.CreateBook)
	router.POST("/book/update", handler.UpdateBook)
}
