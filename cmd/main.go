package main

import (
	"bookstore/internal/middleware"
	"bookstore/internal/router"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	headerLog := "BookStore"

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DB_PORT")

	// connection to db
	conn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		dbname,
		port,
	)

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
