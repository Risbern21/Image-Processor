package server

import (
	"images/routes"

	"github.com/gofiber/fiber/v2"
)

func errorHandler(c *fiber.Ctx, e error) error {
	msg := e.Error()
	return c.Status(fiber.StatusInternalServerError).JSON(msg)
}

var notFoundHandler = func(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).
		JSON("requested resource was not found")
}

func addRoutes(app *fiber.App) {
	baseRouter := app.Group("/")
	// app.Get("/swagger/")

	routes.Images(baseRouter)
	routes.Users(baseRouter)
}
