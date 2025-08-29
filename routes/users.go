package routes

import (
	"images/controllers/users"

	"github.com/gofiber/fiber/v2"
)

func Users(r fiber.Router) {
	userRoutes := r.Group("/users")

	userRoutes.Post("/", users.Create)

	userRoutes.Get("/:id", users.Get)
	userRoutes.Put("/:id", users.Update)
	userRoutes.Delete("/:id", users.Delete)
}
