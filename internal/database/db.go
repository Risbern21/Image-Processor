package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Client() *gorm.DB {
	return db
}

func Connect() {
	dsn := "user=postgres password=postgres database=postgres sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("err while opening database : %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("err before pinging database : %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("unable to ping database : %v", err)
	}
	log.Println("successfully connected to the database")
}
