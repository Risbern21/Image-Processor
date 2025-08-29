package main

import (
	"fmt"
	"images/cmd/app"
	"log"
)

func main() {
	app.SetUp()

	app := app.New()

	fmt.Print("app listening on port 42069")
	log.Fatal(app.Listen(":42069"))
}
