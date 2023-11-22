package db

import (
	"log"
	"os"

	models "darkness-awakens/db/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("Failed to connect to database")
	}
	db.AutoMigrate(&models.User{})
	return db
}
