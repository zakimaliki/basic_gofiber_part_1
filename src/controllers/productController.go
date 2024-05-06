package controllers

import (
	"fmt"
	"gofiber/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var products = []models.Product{
	{ID: 1, Name: "Product A", Price: 10.99, Stock: 100},
	{ID: 2, Name: "Product B", Price: 20.50, Stock: 50},
	{ID: 3, Name: "Product C", Price: 15.75, Stock: 75},
}

func GetAllProduct(c *fiber.Ctx) error {

	// Kirim data produk dalam format JSON
	return c.JSON(products)
}

func GetDetailProduct(c *fiber.Ctx) error {
	// Dapatkan ID produk dari parameter route
	paramId := c.Params("id")
	id, _ := strconv.Atoi(paramId)

	// Panggil fungsi initProducts untuk mendapatkan data produk
	// products := initProducts()

	// Temukan produk dengan ID yang sesuai
	var foundProduct models.Product
	for _, p := range products {
		if p.ID == id {
			foundProduct = p
			break
		}
	}

	// Kirim detail produk dalam format JSON
	return c.JSON(foundProduct)
}

func CreateProduct(c *fiber.Ctx) error {
	// Parse data yang dikirim oleh klien
	var newProduct models.Product
	if err := c.BodyParser(&newProduct); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
		return err
	}

	// Panggil fungsi initProducts untuk mendapatkan data produk
	// products := initProducts()

	// Generate ID untuk produk baru (misalnya, ID terakhir + 1)
	newProduct.ID = len(products) + 1

	// Tambahkan produk baru ke daftar produk
	products = append(products, newProduct)

	// Kirim respons
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"product": newProduct,
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	// Dapatkan ID produk dari parameter route
	id, _ := strconv.Atoi(c.Params("id"))

	// Parse data yang dikirim oleh klien
	var updatedProduct models.Product
	if err := c.BodyParser(&updatedProduct); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
		return err
	}

	// Temukan produk dengan ID yang sesuai
	var foundIndex int = -1
	for i, p := range products {
		if p.ID == id {
			foundIndex = i
			break
		}
	}

	// Jika produk ditemukan, perbarui produk
	if foundIndex != -1 {
		products[foundIndex] = updatedProduct
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": fmt.Sprintf("Product with ID %d updated successfully", id),
			"product": updatedProduct,
		})
	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Product with ID %d not found", id),
		})
	}
}

func DeleteProduct(c *fiber.Ctx) error {
	// Dapatkan ID produk dari parameter route
	id, _ := strconv.Atoi(c.Params("id"))

	// Temukan indeks produk yang akan dihapus
	var foundIndex int = -1
	for i, p := range products {
		if p.ID == id {
			foundIndex = i
			break
		}
	}

	// Jika produk ditemukan, hapus produk
	if foundIndex != -1 {
		// Hapus produk dari slice menggunakan teknik slice trick
		products = append(products[:foundIndex], products[foundIndex+1:]...)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": fmt.Sprintf("Product with ID %d deleted successfully", id),
		})
	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("Product with ID %d not found", id),
		})
	}
}
