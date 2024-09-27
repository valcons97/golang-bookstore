package db

import (
	"bookstore/pkg/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
    dbUrl := "postgres://pg:pass@localhost:5432/crud"

    db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

    if err != nil {
        log.Fatalln(err)
    }

    db.AutoMigrate(&models.Book{})

    return db
}
