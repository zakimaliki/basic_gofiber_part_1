package helpers

import (
	"gofiber/src/configs"
	"gofiber/src/models"
)

func Migration() {
	configs.DB.AutoMigrate(&models.Product{}, &models.Category{})
}
