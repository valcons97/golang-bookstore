package main

import (
	"bookstore/internal/migration"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logHeader := "SeedingPhase"

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DB_PORT")

	// Create the connection string
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

	if err := migration.AddUniqueConstraintIfNotExists(sqlDB); err != nil {
		log.Fatalf("[%v] Could not migrate books: %v", logHeader, err)
	}

	log.Printf("[%v] Database seeded with initial books.", logHeader)
}
