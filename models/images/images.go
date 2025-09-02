package images

import (
	"context"
	"images/internal/database"
	"images/internal/dto"
	"images/models/users"
	"images/transformations"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	UserID uuid.UUID `gorm:"not null;type:uuid" json:"user_id"`
	URL    string    `gorm:"not null"           json:"url"`

	User users.Users `gorm:"foreignKey:UserID;references:ID" json:"-"`

	Image  *dto.Image  `gorm:"-" json:"-"`
	Images *dto.Images `gorm:"-" json:"-"`
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
	if err := database.Client().Where("user_id=?", i.UserID).First(i.Image, i.ID); err != nil {
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

type Filters struct {
	Grayscale bool `json:"grayscale"`
	Sepia     bool `json:"sepia"`
	Invert    bool `json:"invert"`
}

type Crop struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X1     int `json:"x1"`
	Y1     int `json:"y1"`
	X2     int `json:"x2"`
	Y2     int `json:"y2"`
}

type Resize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Transformations struct {
	Resize  *Resize  `json:"resize"`
	Crop    *Crop    `json:"crop"`
	Rotate  *float64 `json:"rotate"`
	Format  *string  `json:"format"`
	Filters *Filters `json:"filters"`
	Flip    *string  `json:"flip"`
}

func NewTransformation() *Transformations {
	return &Transformations{}
}

func (t *Transformations) Transform(c context.Context, url string) error {
	// check resize transformations
	if t.Resize != nil {
		if err := transformations.Resize(t.Resize.Width, t.Resize.Height, url); err != nil {
			return err
		}
	}

	// check crop transformations
	if t.Crop != nil {
		if err := transformations.Crop(t.Crop.X1, t.Crop.Y1, t.Crop.X2, t.Crop.Y2, url); err != nil {
			return err
		}
	}

	// check rotate Transformations
	if t.Rotate != nil {
		if err := transformations.Rotate(*t.Rotate, url); err != nil {
			return err
		}
	}

	// check format transformations
	if t.Format != nil {
		if err := transformations.ConvertFormat(url, *t.Format); err != nil {
			return err
		}
	}

	// check filters transformations
	if t.Filters != nil {
		if t.Filters.Grayscale {
			if err := transformations.Filters("grayscale", url); err != nil {
				return err
			}
		}
		if t.Filters.Invert {
			if err := transformations.Filters("invert", url); err != nil {
				return err
			}
		}
		if t.Filters.Sepia {
			if err := transformations.Filters("sepia", url); err != nil {
				return err
			}
		}
	}

	// check flip transformations
	if t.Flip != nil {
		if err := transformations.Flip(*t.Flip, url); err != nil {
			return err
		}
	}

	// check

	return nil
}
