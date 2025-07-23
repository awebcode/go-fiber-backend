package routes

import (
	"fiber-app/controllers"
 "fiber-app/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Get("/profile", middleware.IsAuthenticated, controllers.Profile)
	app.Post("/logout", controllers.Logout)
}
