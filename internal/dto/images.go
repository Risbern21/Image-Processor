package dto

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID        uint       `json:"id"`
	UserID    uuid.UUID  `json:"user_id:"`
	URL       string     `json:"url"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Images struct {
	Images []Image `json:"images" gorm:"-"`
}
