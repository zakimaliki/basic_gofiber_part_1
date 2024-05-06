package routes

import (
	"gofiber/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	app.Get("/products", controllers.GetAllProduct)
	app.Get("/product/:id", controllers.GetDetailProduct)
	app.Post("/product", controllers.CreateProduct)
	app.Put("/product/:id", controllers.UpdateProduct)
	app.Delete("/product/:id", controllers.DeleteProduct)
}
