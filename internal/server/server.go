package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var app *fiber.App

func New() *fiber.App {
	return app
}

func SetUp() {
	app = fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
		BodyLimit:    16 * 1024 * 1024,
	})

	app.Use(logger.New())
	defer app.Use(notFoundHandler)
	defer app.Use(recover.New())

	addRoutes(app)
}
