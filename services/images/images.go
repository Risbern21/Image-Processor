package images

import (
	"context"
	"images/internal/dto"
	"images/models/images"

	"github.com/google/uuid"
)

type Image struct {
	ID     uint
	UserID uuid.UUID

	Image  *dto.Image
	Images *dto.Images
}

func New() *Image {
	return &Image{}
}

func (i *Image) Create(c context.Context) error {
	m := images.New()
	m.Image = i.Image
	m.UserID = i.UserID

	if err := m.Create(c); err != nil {
		return err
	}
	i.ID = m.ID
	return nil
}

func (i *Image) Get(c context.Context) error {
	m := images.New()
	m.Images = i.Images
	m.ID = i.ID
	m.UserID = i.UserID
	if err := m.Get(c); err != nil {
		return err
	}
	return nil
}

func (i *Image) GetByID(c context.Context) error {
	m := images.New()
	m.Image = i.Image
	m.ID = i.ID
	m.UserID = i.UserID

	if err := m.GetByID(c); err != nil {
		return err
	}
	return nil
}

func (i *Image) Update(c context.Context) error {
	m := images.New()
	m.Image = i.Image
	m.ID = i.ID
	m.UserID = i.UserID

	if err := m.Update(c); err != nil {
		return err
	}
	return nil
}

func (i *Image) Delete(c context.Context) error {
	m := images.New()
	m.ID = i.ID
	m.UserID = i.UserID

	if err := m.Delete(c); err != nil {
		return err
	}
	return nil
}
