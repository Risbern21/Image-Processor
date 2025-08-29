package routes

import "github.com/gofiber/fiber/v2"

func Images(r fiber.Router) {
	images := r.Group("/users/:id/images")

	images.Post("/", nil)
	images.Get("/", nil)

	images.Delete("/:i_id", nil)
	images.Get("/:i_id", nil)
	images.Put("/:i_id", nil)
}
