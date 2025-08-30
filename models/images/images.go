package images

import (
	"images/internal/database"
	"images/internal/dto"
	"images/models/users"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	UserID uuid.UUID `gorm:"not null;type:uuid" json:"user_id"`

	User users.Users `gorm:"foreignKey:UserID;references:ID" json:"-"`

	Image  *dto.Image  `gorm:"-"`
	Images *dto.Images `gorm:"-"`
}

func New() *Image {
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

func (i *Image) Get(c *fiber.Ctx) error {
	if err := database.Client().Table("images").Where("user_id", i.UserID).Find(&i.Images.Images); err != nil {
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
	if err := database.Client().Where("user_id", i.UserID).Unscoped().Delete(&i, i.ID); err != nil {
		return err.Error
	}
	return nil
}
