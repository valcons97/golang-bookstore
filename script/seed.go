package main

import (
	"bookstore/pkg/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// this seeding should only called once when docker-compose
func main() {
	logHeader := "SeedingPhase"
	// connection to db
	conn := "host=db user=user password=password dbname=bookstore port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})

	if err != nil {
		log.Fatalf("[%v] Seeding the initial db failed: %v", logHeader, err)
	}

	if err := models.Migrate(db); err != nil {
		log.Fatalf("[%v] Could not migrate books: %v", logHeader, err)
	}

	// Seed the database with random books

	models.SeedBooks(db, 10)

	log.Printf("[%v] Database seeded with initial books.", logHeader)
}
