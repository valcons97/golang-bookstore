package main

import (
	"bookstore/internal/migration"
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

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[%v] Failed to retrieve sql.DB from GORM: %v", logHeader, err)
	}

	if err := migration.Migrate(sqlDB); err != nil {
		log.Fatalf("[%v] Could not migrate books: %v", logHeader, err)
	}

	// Seed the database with books
	migration.SeedBooks(sqlDB)

	log.Printf("[%v] Database seeded with initial books.", logHeader)
}
