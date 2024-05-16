package controllers

import (
	"gofiber/src/helpers"
	"gofiber/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
)

func GetAllCategories(c *fiber.Ctx) error {
	keyword := c.Query("search")
	categories := models.SelectAllCategories(keyword)
	return c.JSON(categories)
}

func GetCategoryById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	category := models.SelectCategoryById(id)
	return c.JSON(category)
}

func CreateCategory(c *fiber.Ctx) error {
	var category map[string]interface{}
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	category = helpers.XSSMiddleware(category)

	// Convert map to Category model using mapstructure
	var newCategory models.Category
	mapstructure.Decode(category, &newCategory)

	errors := helpers.ValidateStruct(newCategory)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	models.PostCategory(&newCategory)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Category created successfully",
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var updatedCategory map[string]interface{}
	if err := c.BodyParser(&updatedCategory); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
		return err
	}

	updatedCategory = helpers.XSSMiddleware(updatedCategory)

	// Convert map to Category model using mapstructure
	var newUpdatedCategory models.Category
	mapstructure.Decode(updatedCategory, &newUpdatedCategory)

	errors := helpers.ValidateStruct(updatedCategory)

	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	models.UpdateCategory(id, &newUpdatedCategory)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category updated successfully",
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	models.DeleteCategory(id)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}
