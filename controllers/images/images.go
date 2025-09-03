package images

import (
	"errors"
	"fmt"
	"images/internal/dto"
	"images/models/images"
	"images/models/users"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Upload(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	i := images.New()

	if form, err := ctx.MultipartForm(); err == nil {
		file := form.File["file"][0]

		if err := ctx.SaveFile(file, fmt.Sprintf("./assets/%s", file.Filename)); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).
				JSON("unable to save the image")
		}

		i.URL = fmt.Sprintf("./assets/%s", file.Filename)
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON("no image file received")
	}

	uID := ctx.Params("id")
	userID, err := uuid.Parse(uID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}
	// check  if user exists
	u := users.New()
	u.ID = userID

	if err := u.Get(c); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON("user not found")
	}

	i.UserID = userID

	if err := i.Create(c); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ctx.Status(fiber.StatusBadRequest).
				JSON("image already exists")
		}
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return ctx.Status(fiber.StatusCreated).JSON(i)
}

func Get(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	uID := ctx.Params("id")
	userID, err := uuid.Parse(uID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m := images.New()
	m.UserID = userID
	m.Images = &dto.Images{}

	if err := m.Get(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON("requested resources were not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return ctx.Status(fiber.StatusOK).JSON(m.Images.Images)
}

func GetByID(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	iID := ctx.Params("i_id")
	imageID, err := strconv.Atoi(iID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid image id")
	}

	uID := ctx.Params("id")
	userID, err := uuid.Parse(uID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	m := images.New()
	m.ID = uint(imageID)
	m.UserID = userID
	m.Image = &dto.Image{}

	if err := m.GetByID(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON("requested resources were not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return ctx.Status(fiber.StatusOK).JSON(m.Image)
}

func Delete(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	iID := ctx.Params("i_id")
	imageID, err := strconv.Atoi(iID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid image id")
	}

	uID := ctx.Params("id")
	userID, err := uuid.Parse(uID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	u := users.New()
	u.ID = userID

	if err := u.Get(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON("user not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	m := images.New()
	m.ID = uint(imageID)
	m.UserID = userID

	if err := m.Delete(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON("image not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func Transform(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	imageID, err := strconv.Atoi(ctx.Params("i_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid image id")
	}

	userID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	// check  if user exists
	u := users.New()
	u.ID = userID

	if err := u.Get(c); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON("user not found")
	}

	// check if image exists
	i := images.New()
	i.ID = uint(imageID)
	i.UserID = userID
	i.Image = &dto.Image{}

	if err := i.GetByID(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON("image not found or doesnt exist")
		}
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	transformations := images.NewTransformation()
	if err := ctx.BodyParser(transformations); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid input body")
	}

	destURL, err := transformations.Transform(c, i.Image.URL)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("something went wrong")
	}

	destImg := images.New()
	destImg.UserID = userID
	destImg.URL = "./dest/" + destURL

	if err := destImg.Create(c); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ctx.Status(fiber.StatusBadRequest).
				JSON("image already exists")
		}
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("internal server error unable to create destination image")
	}

	return ctx.Status(fiber.StatusOK).JSON(destImg)
}

func Download(ctx *fiber.Ctx) error {
	c := ctx.UserContext()

	imageID, err := strconv.Atoi(ctx.Params("i_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid image id")
	}
	userID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("invalid user id")
	}

	u := users.New()
	u.ID = userID
	if err := u.Get(c); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON("user not found")
	}

	i := images.New()
	i.ID = uint(imageID)
	i.UserID = userID
	i.Image = &dto.Image{}
	if err := i.GetByID(c); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON("image not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("internal server error")
	}

	filePath := i.Image.URL
	fmt.Println(filePath)

	f, err := os.Open(filePath)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON("file not found")
	}
	defer f.Close()

	// read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON("unable to detect file type")
	}
	contentType := http.DetectContentType(buffer)

	f.Seek(0, io.SeekStart)

	fileName := filepath.Base(filePath)

	// set headers to send file
	ctx.Response().Header.Set(
		"Content-Disposition",
		"attachment; filename="+fileName,
	)
	ctx.Response().Header.Set("Content-Type", contentType)
	return ctx.SendFile(filePath)
}
