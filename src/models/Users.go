package models

import (
	"gofiber/src/configs"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
}

func PostUser(item *User) error {
	result := configs.DB.Create(&item)
	return result.Error
}

func FindEmail(input *User) []User {
	items := []User{}
	configs.DB.Raw("SELECT * FROM users WHERE email = ?", input.Email).Scan(&items)
	return items
}
