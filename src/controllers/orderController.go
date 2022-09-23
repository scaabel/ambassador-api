package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/gofiber/fiber/v2"
)

func Orders(c *fiber.Ctx) error {
	var orders []models.Order

	db, _ := database.GetDatabaseConnection()

	db.Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Total = order.GetTotal()
	}

	return c.JSON(&orders)
}
