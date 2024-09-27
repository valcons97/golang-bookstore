package main

import (
	"bookstore/pkg/models"
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
        log.Fatalf("[%v]Could not connect to the database: %v",headerLog, err)
    }

    router := gin.Default()

    router.GET("/books", func(c *gin.Context) {
		var books []models.Book
		db.Find(&books)
		c.JSON(http.StatusOK, books)
	})

	router.POST("/books", func(c *gin.Context) {
		var book models.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&book)
		c.JSON(http.StatusCreated, book)
	})

    if err := router.Run(":8080"); err != nil {
		log.Fatalf("[%v]Could not run server: %v",headerLog, err)
	}
}
