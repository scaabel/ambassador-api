package controllers

import (
	"ambassador/src/database"
	"ambassador/src/middlewares"
	"ambassador/src/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
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
		IsAmbassador: strings.Contains(ctx.Path(), "/api/ambassador"),
	}

	user.SetPassword(data["password"])

	db, _ := database.GetDatabaseConnection()

	db.Create(&user)

	return ctx.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user *models.User

	db, _ := database.GetDatabaseConnection()

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

	isAmbassador := strings.Contains(c.Path(), "/api/ambassador")

	var scope string

	scope = "admin"

	if isAmbassador {
		scope = "ambassador"
	}

	if !isAmbassador && user.IsAmbassador {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized!",
		})
	}

	token, err := middlewares.GenerateJWT(user.Id, scope)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid credentials!",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"message": "Success login"})
}

func User(c *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(c)

	db, conErr := database.GetDatabaseConnection()

	if conErr != nil {
		return c.JSON(fiber.StatusInternalServerError)
	}

	var user *models.User

	db.Where("id = ?", id).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"message": "Success logout"})
}

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := middlewares.GetUserId(c)

	db, _ := database.GetDatabaseConnection()

	user := models.User{
		Name:  data["name"],
		Email: data["email"],
	}

	user.Id = id

	db.Model(&user).Updates(&user)

	return c.JSON(&user)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)

		return c.JSON(fiber.Map{
			"message": "Password does not match!",
		})
	}

	id, _ := middlewares.GetUserId(c)

	db, _ := database.GetDatabaseConnection()

	user := models.User{}

	user.Id = id

	user.SetPassword(data["password"])

	db.Model(&user).Updates(&user)

	return c.JSON(&user)
}
