package users

import (
	"context"
	"images/internal/database"

	"github.com/google/uuid"
)

type Users struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Username string    `gorm:"not null"                             json:"username"`
	Email    string    `gorm:"unique;not null"                      json:"email"`
	Password string    `                                            json:"password"`
}

func New() *Users {
	return &Users{}
}

func (u *Users) Create(c context.Context) error {
	if err := database.Client().Create(u); err != nil {
		return err.Error
	}
	return nil
}

func (u *Users) Get(c context.Context) error {
	if err := database.Client().First(&u, u.ID); err != nil {
		return err.Error
	}
	return nil
}

func (u *Users) Update(c context.Context) error {
	if err := database.Client().Save(u); err != nil {
		return err.Error
	}
	return nil
}

func (u *Users) Delete(c context.Context) error {
	if err := database.Client().Delete(u, u.ID); err != nil {
		return err.Error
	}
	return nil
}
