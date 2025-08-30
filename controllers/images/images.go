package images

import (
	"errors"
	"fmt"
	"images/internal/dto"
	"images/models/images"
	"images/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Create(c *fiber.Ctx) error {
	if form, err := c.MultipartForm(); err == nil {
		file := form.File["file"][0]

		if err := c.SaveFile(file, fmt.Sprintf("./assets/%s", file.Filename)); err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON("unable to save the image")
		}
	}

	uID := c.Params("id")
	userID, err := uuid.Parse(uID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m := images.New()
	m.UserID = userID

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

	m := images.New()
	m.UserID = userID
	m.Images = &dto.Images{}

	if err := m.Get(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).
				JSON("requested resources were not found")
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	path, err := utils.ConvertFormat("assets/image.jpeg", "png")
	if err != nil {
		fmt.Print("errror inside\n")
		return c.Status(fiber.StatusInternalServerError).
			JSON("error occurred while converting the formats")
	}

	fmt.Println(path)

	return c.Status(fiber.StatusOK).JSON(m.Images.Images)
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

	m := images.New()
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

	m := images.New()
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

	m := images.New()
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
