package images

import (
	"context"
	"fmt"
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
	fmt.Println("helllo")
	fmt.Println(i)
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
	Mirror  *string  `json:"mirror"`
}

func NewTransformation() *Transformations {
	return &Transformations{}
}

func (t *Transformations) Transform(
	c context.Context,
	url string,
) (string, error) {
	var destUrl string = ""
	var err error

	if t.Resize != nil {
		if destUrl, err = transformations.Resize(t.Resize.Width, t.Resize.Height, url); err != nil {
			return "", err
		}
	}

	if t.Crop != nil {
		if t.Crop.X2-t.Crop.X1 < 10 || t.Crop.Y2-t.Crop.Y1 < 10 {
			return "", fmt.Errorf(
				"cannot crop image into such small proportions",
			)
		}
		if destUrl, err = transformations.Crop(t.Crop.X1, t.Crop.Y1, t.Crop.X2, t.Crop.Y2, url); err != nil {
			return "", err
		}
	}

	if t.Rotate != nil {
		if destUrl, err = transformations.Rotate(*t.Rotate, url); err != nil {
			return "", err
		}
	}

	if t.Format != nil {
		if destUrl, err = transformations.ConvertFormat(url, *t.Format); err != nil {
			return "", err
		}
	}

	if t.Filters != nil {
		if t.Filters.Grayscale {
			if destUrl, err = transformations.Filters("grayscale", url); err != nil {
				return "", err
			}
		}
		if t.Filters.Invert {
			if destUrl, err = transformations.Filters("invert", url); err != nil {
				return "", err
			}
		}
		if t.Filters.Sepia {
			if destUrl, err = transformations.Filters("sepia", url); err != nil {
				return "", err
			}
		}
	}

	if t.Flip != nil {
		if destUrl, err = transformations.Flip(*t.Flip, url); err != nil {
			return "", err
		}
	}

	if t.Mirror != nil {
		if destUrl, err = transformations.Mirror(url, *t.Mirror); err != nil {
			return "", err
		}
	}

	return destUrl, nil
}
