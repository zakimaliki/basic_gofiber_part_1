package routes

import (
	"gofiber/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	// Product routes
	app.Get("/products", controllers.GetAllProducts)
	app.Get("/product/:id", controllers.GetProductById)
	app.Post("/product", controllers.CreateProduct)
	app.Put("/product/:id", controllers.UpdateProduct)
	app.Delete("/product/:id", controllers.DeleteProduct)

	// Category routes
	app.Get("/categories", controllers.GetAllCategories)
	// app.Get("/categories", middlewares.JwtMiddleware(), controllers.GetAllCategories)

	app.Get("/category/:id", controllers.GetCategoryById)
	app.Post("/category", controllers.CreateCategory)
	app.Put("/category/:id", controllers.UpdateCategory)
	app.Delete("/category/:id", controllers.DeleteCategory)

	// User Routes
	app.Post("/register", controllers.RegisterUser)
	app.Post("/login", controllers.LoginUser)
}
