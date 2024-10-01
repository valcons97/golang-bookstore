package router

import (
	"bookstore/internal/handler"
	"bookstore/internal/repository"
	"bookstore/internal/service"

	"database/sql"

	"github.com/gin-gonic/gin"
)

func CustomerRouter(router *gin.Engine, db *sql.DB) {
	repo := repository.NewCustomerRepository(db)
	svc := service.NewCustomerService(repo)
	handler := handler.NewCustomerHandler(svc)

	// Define the routes
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

}
