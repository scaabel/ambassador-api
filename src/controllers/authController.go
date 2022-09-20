package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	var data map[string]string

	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		ctx.Status(400)

		return ctx.JSON(fiber.Map{
			"message": "Password does not match!",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)

	user := models.User{
		Name:         data["name"],
		Email:        data["email"],
		Password:     password,
		IsAmbassador: false,
	}

	db, conErr := database.GetDatabaseConnection()

	if conErr != nil {
		ctx.Status(500)
		return ctx.JSON(fiber.Map{
			"message": "Service unavailable!",
		})
	}

	db.Create(&user)

	return ctx.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	db, conErr := database.GetDatabaseConnection()

	if conErr != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Service unavailable!",
		})
	}

	db.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid email or password!",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid email or password!",
		})
	}

	return c.JSON(user)
}
