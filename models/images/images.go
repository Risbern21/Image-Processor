package images

import (
	"images/internal/database"
	"images/models/users"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	UserID uuid.UUID `gorm:"not null;type:uuid" json:"user_id"`

	User users.Users `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func NewImage() *Image {
	return &Image{}
}

func (i *Image) Create(c *fiber.Ctx) error {
	if err := database.Client().Create(&i); err != nil {
		return err.Error
	}
	return nil
}

func (i *Image) GetByID(c *fiber.Ctx) error {
	if err := database.Client().Where("user_id=?", i.UserID).First(i, i.ID); err != nil {
		return err.Error
	}
	return nil
}

func (i *Image) Update(c *fiber.Ctx) error {
	if err := database.Client().Where("user_id", i.UserID).Save(i); err != nil {
		return err.Error
	}
	return nil
}

func (i *Image) Delete(c *fiber.Ctx) error {
	if err := database.Client().Where("user_id", i.UserID).Delete(i, i.ID); err != nil {
		return err.Error
	}
	return nil
}

type Images struct {
	UserID uuid.UUID `json:"user_id"`
	Images []Image   `json:"images"`

	User users.Users `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func NewImages() *Images {
	return &Images{}
}

func (i *Images) Get(c *fiber.Ctx) error {
	if err := database.Client().Where("user_id", i.UserID).Find(i.Images); err != nil {
		return err.Error
	}
	return nil
}
