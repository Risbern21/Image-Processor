package users

import (
	"errors"
	"images/internal/dto"
	users "images/models/users"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Create(c *fiber.Ctx) error {
	m := users.New()

	var user dto.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid input body")
	}

	m.Username = user.Username
	m.Email = user.Email

	if err := m.Create(c); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.Status(fiber.StatusBadRequest).
				JSON("this email is already registered")
		}

		return c.Status(fiber.StatusInternalServerError).
			JSON("something went wrong")
	}

	return c.Status(fiber.StatusCreated).JSON(m)
}

func Get(c *fiber.Ctx) error {
	m := users.New()

	id := c.Params("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m.ID = userID

	if err := m.Get(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("user not found")
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(m)
}

func Update(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m := users.New()
	var user dto.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid input body")
	}

	m.ID = userID
	m.Username = user.Username
	m.Email = user.Email

	if err := m.Update(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("user not found")
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m := users.New()

	m.ID = userID
	if err := m.Delete(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("user not found")
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
