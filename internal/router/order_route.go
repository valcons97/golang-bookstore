package router

import (
	"bookstore/internal/handler"
	"bookstore/internal/repository"
	"bookstore/internal/service"

	"database/sql"

	"github.com/gin-gonic/gin"
)

func OrderRouter(router *gin.Engine, db *sql.DB, authMiddleware gin.HandlerFunc) {
	repo := repository.NewOrderRepository(db) // Create the repository
	svc := service.NewOrderService(repo)      // Create the service
	handler := handler.NewOrderHandler(svc)   // Create the handler

	orderRoutes := router.Group("/orders", authMiddleware)

	// Define the routes
	orderRoutes.POST("/add", handler.AddBookToOrder)
}
