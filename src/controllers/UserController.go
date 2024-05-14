package controllers

import (
	"gofiber/src/helpers"
	"gofiber/src/models"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashPassword)

	models.PostUser(&user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Register successfully",
	})
}

func LoginUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	validateEmail := models.FindEmail(&user)
	if len(validateEmail) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Email is not Found",
		})
	}

	var PasswordSecond string
	for _, user := range validateEmail {
		PasswordSecond = user.Password
	}

	if err := bcrypt.CompareHashAndPassword([]byte(PasswordSecond), []byte(user.Password)); err != nil || user.Password == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Password invalid",
		})
	}

	jwtKey := os.Getenv("SECRETKEY")
	token, _ := helpers.GenerateToken(jwtKey, user.Email)
	item := map[string]string{
		"Email": user.Email,
		"Token": token,
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Login successfully",
		"data":    item,
	})
}
