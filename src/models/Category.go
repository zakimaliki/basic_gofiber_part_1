package models

import (
	"gofiber/src/configs"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name     string       `json:"name" validate:"required,min=3,max=50"`
	Image    string       `json:"image"`
	Products []ApiProduct `json:"products"`
}

type ApiProduct struct {
	Name       string  `json:"name" `
	Price      float64 `json:"price" `
	Stock      int     `json:"stock" `
	CategoryID uint    `json:"category_id"`
}

func SelectAllCategories(sort, name string) []*Category {
	var categories []*Category
	name = "%" + name + "%"
	configs.DB.Preload("Products", func(db *gorm.DB) *gorm.DB {
		var items []*ApiProduct
		return db.Model(&Product{}).Find(&items)
	}).Order(sort).Where("name LIKE ?", name).Find(&categories)
	return categories
}

func SelectCategoryById(id int) *Category {
	var category Category
	configs.DB.Preload("Products",
		func(db *gorm.DB) *gorm.DB {
			var items []*ApiProduct
			return db.Model(&Product{}).Find(&items)
		}).First(&category, "id = ?", id)
	return &category
}

func PostCategory(category *Category) error {
	result := configs.DB.Create(&category)
	return result.Error
}

func UpdateCategory(id int, updatedCategory *Category) error {
	result := configs.DB.Model(&Category{}).Where("id = ?", id).Updates(updatedCategory)
	return result.Error
}

func DeleteCategory(id int) error {
	result := configs.DB.Delete(&Category{}, "id = ?", id)
	return result.Error
}
