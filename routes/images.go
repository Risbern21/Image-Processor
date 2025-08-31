package routes

import (
	"images/controllers/images"

	"github.com/gofiber/fiber/v2"
)

func Images(r fiber.Router) {
	imageRoutes := r.Group("/users/:id/images")

	imageRoutes.Post("/", images.Upload)
	imageRoutes.Get("/", images.Get)

	imageRoutes.Get("/:i_id", images.GetByID)
	imageRoutes.Put("/:i_id", images.Edit)
	imageRoutes.Delete("/:i_id", images.Delete)

	imageRoutes.Post("/:i_id/transform", images.Transform)
}
