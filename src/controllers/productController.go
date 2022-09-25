package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Products(c *fiber.Ctx) error {
	var products []models.Product

	db, _ := database.GetDatabaseConnection()

	db.Find(&products)

	return c.JSON(products)
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	db, _ := database.GetDatabaseConnection()

	db.Create(&product)

	return c.JSON(product)
}

func GetProduct(c *fiber.Ctx) error {
	var product models.Product

	id, _ := strconv.Atoi(c.Params("id"))

	product.Id = uint(id)

	db, _ := database.GetDatabaseConnection()

	db.Find(&product)

	return c.JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}

	product.Id = uint(id)

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	db, _ := database.GetDatabaseConnection()

	db.Model(&product).Updates(&product)

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}

	product.Id = uint(id)

	db, _ := database.GetDatabaseConnection()

	db.Delete(&product)

	return c.JSON(fiber.Map{"message": "Success"})
}

func ProudctsFrontend(c *fiber.Ctx) error {

	var products []models.Product
	var ctx = context.Background()

	result, err := database.Cache.Get(ctx, "products_frontend").Result()

	db, _ := database.GetDatabaseConnection()

	if err != nil {
		db.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		if errKey := database.Cache.Set(ctx, "products_frontend", bytes, 30*time.Minute).Err(); errKey != nil {
			panic(errKey)
		}
	}

	if result != "" {
		json.Unmarshal([]byte(result), &products)
	}

	return c.JSON(products)
}

func ProudctsBackend(c *fiber.Ctx) error {

	var products []models.Product
	var ctx = context.Background()

	result, err := database.Cache.Get(ctx, "products_backend").Result()

	db, _ := database.GetDatabaseConnection()

	if err != nil {
		db.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		if errKey := database.Cache.Set(ctx, "products_backend", bytes, 30*time.Minute).Err(); errKey != nil {
			panic(errKey)
		}
	}

	if result != "" {
		json.Unmarshal([]byte(result), &products)
	}

	var searchedProducts []models.Product

	if search := c.Query("search"); search != "" {
		lowerCase := strings.ToLower(search)

		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), lowerCase) || strings.Contains(strings.ToLower(product.Description), lowerCase) {
				searchedProducts = append(searchedProducts, product)
			}
		}

		return c.JSON(searchedProducts)
	}

	return c.JSON(products)
}
