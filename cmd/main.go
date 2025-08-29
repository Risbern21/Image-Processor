package main

import (
	"fmt"
	"images/cmd/app"
	"images/internal/database"
	"log"
)

func main() {
	app.SetUp()
	database.Connect()

	app := app.New()

	fmt.Print("app listening on port 42069")
	log.Fatal(app.Listen(":42069"))
}
