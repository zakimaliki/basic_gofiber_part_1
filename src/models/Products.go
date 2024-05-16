package models

import (
	"gofiber/src/configs"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name       string   `json:"name" validate:"required,min=3,max=100"`
	Price      float64  `json:"price" validate:"required,min=0"`
	Stock      int      `json:"stock" validate:"required,min=0"`
	CategoryID uint     `json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID"`
}

func SelectAllProducts(sort, name string, limit, offset int) []*Product {
	var items []*Product
	name = "%" + name + "%"
	configs.DB.Preload("Category").Order(sort).Limit(limit).Offset(offset).Where("name LIKE ?", name).Find(&items)
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

func CountData() int64 {
	var result int64
	configs.DB.Table("products").Count(&result)
	return result
}
