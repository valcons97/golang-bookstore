package router

import (
	"bookstore/internal/handler"
	"bookstore/internal/repository"
	"bookstore/internal/service"

	"database/sql"

	"github.com/gin-gonic/gin"
)

func OrderRouter(router *gin.Engine, db *sql.DB, authMiddleware gin.HandlerFunc) {
	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo)
	handler := handler.NewOrderHandler(svc)

	orderRoutes := router.Group("/orders", authMiddleware)

	// Define the routes
	orderRoutes.POST("/add", handler.AddToCart)
	orderRoutes.POST("/pay", handler.PayOrder)
	orderRoutes.POST("/delete", handler.RemoveFromCart)
	orderRoutes.GET("/cart", handler.GetCart)
	orderRoutes.POST("/history", handler.GetOrderHistory)

}
