package images

import (
	"context"
	"images/internal/database"
	"images/internal/dto"
	"images/models/users"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Filters struct {
	Grayscale bool `json:"grayscale"`
	Sepia     bool `json:"sepia"`
}

type Crop struct {
	Width  uint16 `json:"width"`
	Height uint16 `json:"height"`
	X      int16  `json:"x"`
	Y      int16  `json:"y"`
}

type Resize struct {
	Width  uint16 `json:"width"`
	Height uint16 `json:"height"`
}

type Transformations struct {
	Resize  Resize  `json:"resize"`
	Crop    Crop    `json:"crop"`
	Rotate  int16   `json:"rotate"`
	Format  string  `json:"format"`
	Filters Filters `json:"filters"`
}

type Image struct {
	gorm.Model
	UserID          uuid.UUID       `gorm:"not null;type:uuid" json:"user_id"`
	URL             string          `gorm:"not null"           json:"url"`
	Transformations Transformations `                          json:"transformations"`

	User users.Users `gorm:"foreignKey:UserID;references:ID" json:"-"`

	Image  *dto.Image  `gorm:"-"`
	Images *dto.Images `gorm:"-"`
}

func New() *Image {
	return &Image{}
}

func (i *Image) Create(c context.Context) error {
	if err := database.Client().Create(&i); err != nil {
		return err.Error
	}

	return nil
}

func (i *Image) GetByID(c context.Context) error {
	if err := database.Client().Where("user_id=?", i.UserID).First(i, i.ID); err != nil {
		return err.Error
	}
	return nil
}

func (i *Image) Get(c context.Context) error {
	if err := database.Client().Table("images").Where("user_id", i.UserID).Find(&i.Images.Images); err != nil {
		return err.Error
	}
	return nil
}

func (i *Image) Update(c context.Context) error {
	if err := database.Client().Where("user_id", i.UserID).Save(i); err != nil {
		return err.Error
	}
	return nil
}

func (i *Image) Delete(c context.Context) error {
	if err := database.Client().Where("user_id", i.UserID).Unscoped().Delete(&i, i.ID); err != nil {
		return err.Error
	}
	return nil
}
