package controllers

import (
	"gofiber/src/helpers"
	"gofiber/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllCategories(c *fiber.Ctx) error {
	categories := models.SelectAllCategories()
	return c.JSON(categories)
}

func GetCategoryById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	category := models.SelectCategoryById(id)
	return c.JSON(category)
}

func CreateCategory(c *fiber.Ctx) error {
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	errors := helpers.ValidateStruct(category)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	models.PostCategory(&category)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Category created successfully",
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var updatedCategory models.Category
	if err := c.BodyParser(&updatedCategory); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
		return err
	}

	errors := helpers.ValidateStruct(updatedCategory)

	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	models.UpdateCategory(id, &updatedCategory)

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
