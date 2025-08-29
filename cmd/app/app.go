package app

import (
	"images/internal/config"
	"images/internal/database"
	"images/internal/server"
	"log"
)

func SetUp() {
	database.Connect()
	config.AutoMigrate()
	server.SetUp()

	app := server.New()

	if err := app.Listen(":42069"); err != nil {
		log.Fatalf("%v", err)
	}
}
