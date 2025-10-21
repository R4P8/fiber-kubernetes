package routes

import (
	"example/controllers"

	"github.com/gofiber/fiber/v2"
)

func CategoryRoutes(app *fiber.App) {
	api := app.Group("/api/categories")

	api.Get("/", controllers.GetCategories)
	api.Get("/:id", controllers.GetCategory)
	api.Post("/", controllers.CreateCategory)
	api.Put("/:id", controllers.UpdateCategory)
	api.Delete("/:id", controllers.DeleteCategory)
}
