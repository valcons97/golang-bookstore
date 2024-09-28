package main

import (
	"bookstore/internal/book"
	"log"
	"net/http"

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

	router := gin.Default()

	router.GET("/tables", func(c *gin.Context) {
		// Get the table names from the database schema using GORM's Migrator
		tables := []string{}

		migrator := db.Migrator()

		tables, err := migrator.GetTables()
		if err != nil {
			return
		}

		// Return the table names as JSON response
		c.JSON(http.StatusOK, tables)
	})

	book.Router(router, sqlDB)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("[%v]Could not run server: %v", headerLog, err)
	}
}
