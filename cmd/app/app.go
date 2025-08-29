package app

import "github.com/gofiber/fiber/v2"

var app *fiber.App

func New() *fiber.App {
	return app
}

func SetUp() {
	app = fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON("hello world")
	})
}
