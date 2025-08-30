package dto

import (
	"time"

	"github.com/google/uuid"
)

type ImageDTO struct{}

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

type Image struct {
	ID        uint       `json:"id"`
	UserID    uuid.UUID  `json:"user_id:"`
	URL       string     `json:"url"`
	Resize    *Resize    `json:"resize"`
	Crop      *Crop      `json:"crop"`
	Rotate    *int16     `json:"rotate"`
	Format    *string    `json:"format"`
	Filters   *Filters   `json:"filters"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Images struct {
	Images []Image `json:"images" gorm:"-"`
}
