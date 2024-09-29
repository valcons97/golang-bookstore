package main

import (
	"bookstore/internal/middleware"
	"bookstore/internal/model"
	"bookstore/internal/router"
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

	r := gin.Default()

	r.GET("/tables", func(c *gin.Context) {
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

	r.GET("/customers", func(c *gin.Context) {
		// Declare a slice to hold customers
		var customers []model.Customer

		// Use GORM to find all customers
		err := db.Find(&customers).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the customers as JSON response
		c.JSON(http.StatusOK, customers)
	})

	r.GET("/orders", func(c *gin.Context) {
		// Declare a slice to hold customers
		var orders []model.Order

		// Use GORM to find all customers
		err := db.Find(&orders).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the customers as JSON response
		c.JSON(http.StatusOK, orders)
	})

	r.GET("/ordeDetails", func(c *gin.Context) {
		// Declare a slice to hold customers
		var orderDetails []model.OrderDetail

		// Use GORM to find all customers
		err := db.Find(&orderDetails).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the customers as JSON response
		c.JSON(http.StatusOK, orderDetails)
	})

	authMiddleware := middleware.AuthMiddleware()

	router.BookRouter(r, sqlDB)
	router.CustomerRouter(r, sqlDB)
	router.OrderRouter(r, sqlDB, authMiddleware)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("[%v]Could not run server: %v", headerLog, err)
	}
}
