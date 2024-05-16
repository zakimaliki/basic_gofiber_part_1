package controllers

import (
	"fmt"
	"gofiber/src/helpers"
	"gofiber/src/models"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
)

func GetAllProducts(c *fiber.Ctx) error {
	pageOld := c.Query("page")
	limitOld := c.Query("limit")
	page, _ := strconv.Atoi(pageOld)
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitOld)
	if limit == 0 {
		limit = 5
	}
	offset := (page - 1) * limit
	sort := c.Query("sorting")
	if sort == "" {
		sort = "ASC"
	}
	sortby := c.Query("orderBy")
	if sortby == "" {
		sortby = "name"
	}
	sort = sortby + " " + strings.ToLower(sort)
	keyword := c.Query("search")
	products := models.SelectAllProducts(sort, keyword, limit, offset)
	totalData := models.CountData()
	totalPage := math.Ceil(float64(totalData) / float64(limit))
	result := map[string]interface{}{
		"data":        products,
		"currentPage": page,
		"limit":       limit,
		"totalData":   totalData,
		"totalPage":   totalPage,
	}
	return c.JSON(result)
}

func GetProductById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	foundProduct := models.SelectProductById(id)
	if foundProduct == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
		})
	}
	return c.JSON(foundProduct)
}

func CreateProduct(c *fiber.Ctx) error {
	var Product map[string]interface{}
	if err := c.BodyParser(&Product); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
		return err
	}

	Product = helpers.XSSMiddleware(Product)

	var newProduct models.Product
	mapstructure.Decode(Product, &newProduct)

	errors := helpers.ValidateStruct(newProduct)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}
	models.PostProduct(&newProduct)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var updatedProduct map[string]interface{}
	if err := c.BodyParser(&updatedProduct); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
		return err
	}

	updatedProduct = helpers.XSSMiddleware(updatedProduct)

	var newUpdatedProduct models.Product
	mapstructure.Decode(updatedProduct, &newUpdatedProduct)

	errors := helpers.ValidateStruct(newUpdatedProduct)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	err := models.UpdateProduct(id, &newUpdatedProduct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to update product with ID %d", id),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Product with ID %d updated successfully", id),
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	err := models.DeleteProduct(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to delete product with ID %d", id),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Product with ID %d deleted successfully", id),
	})
}
