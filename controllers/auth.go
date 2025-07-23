package controllers

import (
	"fiber-app/database"
	"fiber-app/models"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

var store = session.New()
var validate = validator.New()

type RegisterInput struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func Register(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := validate.Struct(input); err != nil {
		errs := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errs[e.Field()] = fmt.Sprintf("Validation failed on %s", e.Tag())
		}
		return c.Status(400).JSON(errs)
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	user := models.User{Username: input.Username, Email: input.Email, Password: string(hashedPass)}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email or username already taken"})
	}

	return c.JSON(fiber.Map{"message": "User registered"})
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Validation failed"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Wrong password"})
	}

	sess, _ := store.Get(c)
	sess.Set("userID", user.ID)
	sess.Save()

	return c.JSON(fiber.Map{"message": "Logged in"})
}

func Profile(c *fiber.Ctx) error {
	sess, _ := store.Get(c)
	id := sess.Get("userID")
	if id == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{"id": user.ID, "email": user.Email, "username": user.Username})
}

func Logout(c *fiber.Ctx) error {
	sess, _ := store.Get(c)
	sess.Destroy()
	return c.JSON(fiber.Map{"message": "Logged out"})
}
