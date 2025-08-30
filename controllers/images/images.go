package images

import (
	"errors"
	"images/models/images"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Create(c *fiber.Ctx) error {
	uID := c.Params("id")
	userID, err := uuid.Parse(uID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m := images.NewImage()
	m.UserID = userID

	if err := c.BodyParser(&m); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid input body")
	}

	if err := m.Create(c); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.Status(fiber.StatusBadRequest).JSON("image already exists")
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return c.Status(fiber.StatusCreated).JSON(m)
}

func Get(c *fiber.Ctx) error {
	uID := c.Params("id")
	userID, err := uuid.Parse(uID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m := images.NewImages()
	m.UserID = userID

	if err := m.Get(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).
				JSON("requested resources were not found")
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(m.Images)
}

func GetByID(c *fiber.Ctx) error {
	iID := c.Params("i_id")
	imageID, err := strconv.Atoi(iID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid image id")
	}

	uID := c.Params("id")
	userID, err := uuid.Parse(uID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m := images.NewImage()
	m.ID = uint(imageID)
	m.UserID = userID

	if err := m.GetByID(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).
				JSON("requested resources were not found")
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(m)
}

func Update(c *fiber.Ctx) error {
	iID := c.Params("i_id")
	imageID, err := strconv.Atoi(iID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid image id")
	}

	uID := c.Params("id")
	userID, err := uuid.Parse(uID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m := images.NewImage()
	m.ID = uint(imageID)
	m.UserID = userID

	if err := m.Update(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("image not found")
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func Delete(c *fiber.Ctx) error {
	iID := c.Params("i_id")
	imageID, err := strconv.Atoi(iID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid image id")
	}

	uID := c.Params("id")
	userID, err := uuid.Parse(uID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m := images.NewImage()
	m.ID = uint(imageID)
	m.UserID = userID

	if err := m.Delete(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON("image not found")
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
