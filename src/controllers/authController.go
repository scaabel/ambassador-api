package controllers

import (
	"ambassador/src/config"
	"ambassador/src/database"
	"ambassador/src/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
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

	user := models.User{
		Name:         data["name"],
		Email:        data["email"],
		IsAmbassador: false,
	}

	user.SetPassword(data["password"])

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
		return c.JSON(fiber.StatusInternalServerError)
	}

	db.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid email or password!",
		})
	}

	if err := user.CheckPassword(data["password"]); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid email or password!",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = user.Id
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config("JWT_SECRET")))

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid credentials!",
		})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}
