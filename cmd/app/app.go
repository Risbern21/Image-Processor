package app

import (
	"images/internal/cache"
	"images/internal/config"
	"images/internal/database"
	"images/internal/server"
	"log"
	"os"
)

func SetUp() {
	database.Connect()
	config.AutoMigrate()

	cache.Connect()

	server.SetUp()

	if err := os.MkdirAll("dest", 0755); err != nil {
		log.Fatal(err)
	}

	app := server.New()

	if err := app.Listen(":42069"); err != nil {
		log.Fatalf("%v", err)
	}
}
