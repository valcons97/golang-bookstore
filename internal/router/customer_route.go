package router

import (
	"bookstore/internal/handler"
	"bookstore/internal/repository"
	"bookstore/internal/service"

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
