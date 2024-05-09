package models

import (
	"gofiber/src/configs"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name       string   `json:"name"`
	Price      float64  `json:"price"`
	Stock      int      `json:"stock"`
	CategoryID uint     `json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID"`
}

func SelectAllProducts() []*Product {
	var items []*Product
	configs.DB.Preload("Category").Find(&items)
	return items
}

func SelectProductById(id int) *Product {
	var item Product
	configs.DB.Preload("Category").First(&item, "id = ?", id)
	return &item
}

func PostProduct(item *Product) error {
	result := configs.DB.Create(&item)
	return result.Error
}

func UpdateProduct(id int, newProduct *Product) error {
	var item Product
	result := configs.DB.Model(&item).Where("id = ?", id).Updates(newProduct)
	return result.Error
}

func DeleteProduct(id int) error {
	var item Product
	result := configs.DB.Delete(&item, "id = ?", id)
	return result.Error
}
