package main

import (
	"log"

	"gofiber/src/configs"
	"gofiber/src/helpers"
	"gofiber/src/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*", // Atur "*" untuk mengizinkan akses dari semua origin
		AllowMethods:  "GET,POST,PUT,DELETE",
		AllowHeaders:  "*",
		ExposeHeaders: "Content-Length",
	}))

	configs.InitDB()
	helpers.Migration()
	routes.Router(app)

	app.Listen(":3000")
}
