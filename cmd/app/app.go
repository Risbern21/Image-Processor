package app

import (
	"images/internal/cache"
	"images/internal/config"
	"images/internal/database"
	"images/internal/server"
	"log"
)

func SetUp() {
	database.Connect()
	config.AutoMigrate()

	cache.Connect()

	server.SetUp()

	app := server.New()

	if err := app.Listen(":42069"); err != nil {
		log.Fatalf("%v", err)
	}
}
