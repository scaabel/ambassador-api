package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Link(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var links []models.Link

	db, _ := database.GetDatabaseConnection()

	db.Where("user_id = ?", id).Find(&links)

	for i, link := range links {
		var orders []models.Order

		db.Where("code = ? and complete = true", link.Code).Find(&orders)

		links[i].Orders = orders
	}

	return c.JSON(fiber.Map{"links": links})

}
