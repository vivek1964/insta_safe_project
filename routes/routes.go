package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vivek1964/insta_safe_project/controllers"
)

func Setup(app *fiber.App) {
	app.Post("/transactions", controllers.Transactions)
	app.Get("/statistics", controllers.Statistics)
	app.Delete("/transactions", controllers.Transactions)
}
