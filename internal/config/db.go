package config

import (
	"images/internal/database"
	"images/models/images"
	users "images/models/users"
)

func AutoMigrate() {
	database.Client().AutoMigrate(&users.Users{}, &images.Image{})
}
