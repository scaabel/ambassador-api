package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/gofiber/fiber/v2"
)

func Ambassador(c *fiber.Ctx) error {
	var users []models.User

	db, _ := database.GetDatabaseConnection()

	db.Where("is_ambassador = true").Find(&users)

	return c.JSON(users)
}
