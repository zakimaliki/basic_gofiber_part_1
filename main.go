package main

import (
	"gofiber/src/configs"
	"gofiber/src/helpers"
	"gofiber/src/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()
	configs.InitDB()
	helpers.Migration()
	routes.Router(app)

	app.Listen(":3000")
}
