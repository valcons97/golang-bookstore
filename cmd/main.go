package main

import (
	"bookstore/internal/middleware"
	"bookstore/internal/router"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	headerLog := "BookStore"

	// connection to db
	conn := "host=db user=user password=password dbname=bookstore port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})

	if err != nil {
		log.Fatalf("[%v]Could not connect to the database: %v", headerLog, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve sql.DB from GORM: %v", err)
	}

	r := gin.Default()

	authMiddleware := middleware.AuthMiddleware()

	router.BookRouter(r, sqlDB)
	router.CustomerRouter(r, sqlDB)
	router.OrderRouter(r, sqlDB, authMiddleware)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("[%v]Could not run server: %v", headerLog, err)
	}
}
