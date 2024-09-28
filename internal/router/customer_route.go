package router

import (
	"bookstore/internal/customer/handler"
	"bookstore/internal/customer/repository"
	"bookstore/internal/customer/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func CustomerRouter(router *gin.Engine, db *sql.DB) {
	repo := repository.NewCustomerRepository(db) // Create the repository
	svc := service.NewCustomerService(repo)      // Create the service
	handler := handler.NewCustomerHandler(svc)   // Create the handler

	// Define the routes
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
}
