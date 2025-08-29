package config

import (
	"images/internal/database"
	users "images/models"
)

func AutoMigrate() {
	database.Client().AutoMigrate(&users.Users{})
}
