package controllers

import (
	"ambassador/src/config"
	"ambassador/src/database"
	"ambassador/src/models"
	"strconv"
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

	var user *models.User

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

	payload := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(config.Config("JWT_SECRET")))

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
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Config("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated!",
		})
	}

	payload := token.Claims.(*jwt.StandardClaims)

	db, conErr := database.GetDatabaseConnection()

	if conErr != nil {
		return c.JSON(fiber.StatusInternalServerError)
	}

	var user *models.User

	db.Where("id = ?", payload.Subject).First(&user)

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
