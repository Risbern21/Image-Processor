package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	Username string `json:"username"`
	Email    string `json:"emila"`
	Password string `json:"password"`
}

type User struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
