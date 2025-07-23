package main

import (
	"fiber-app/database"
	"fiber-app/models"
	"fiber-app/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	app := fiber.New()

	app.Use(logger.New())

	database.Connect()
	database.DB.AutoMigrate(&models.User{})

	routes.Setup(app)

	app.Listen(":3000")
}
