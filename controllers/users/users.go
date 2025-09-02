package users

import (
	"errors"
	"images/internal/dto"
	users "images/models/users"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Create(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	m := users.New()

	var user dto.UserDTO

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid input body")
	}

	m.Username = user.Username
	m.Email = user.Email
	m.Password = user.Password

	if err := m.Create(c); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ctx.Status(fiber.StatusBadRequest).
				JSON("this email is already registered")
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON("something went wrong")
	}

	return ctx.Status(fiber.StatusCreated).JSON(m)
}

func Get(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	m := users.New()

	id := ctx.Params("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m.ID = userID

	if err := m.Get(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON("user not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return ctx.Status(fiber.StatusOK).JSON(m)
}

func Update(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	id := ctx.Params("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m := users.New()
	var user dto.UserDTO

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid input body")
	}

	m.ID = userID
	m.Username = user.Username
	m.Email = user.Email
	m.Password = user.Password

	if err := m.Update(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON("user not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func Delete(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	id := ctx.Params("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m := users.New()

	m.ID = userID
	if err := m.Delete(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON("user not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}
