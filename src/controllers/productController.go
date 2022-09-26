package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"encoding/json"
	"sort"
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

	go database.ClearCache("products_backend", "products_frontend")

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

	go database.ClearCache("products_backend", "products_frontend")

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{}

	product.Id = uint(id)

	db, _ := database.GetDatabaseConnection()

	db.Delete(&product)

	go database.ClearCache("products_backend", "products_frontend")

	return c.JSON(fiber.Map{"message": "Success"})
}

func ProudctsFrontend(c *fiber.Ctx) error {
	var products []models.Product

	result, err := database.GetCache("products_frontend")

	db, _ := database.GetDatabaseConnection()

	if err != nil {
		db.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		if errKey := database.SetCache("products_frontend", bytes, 30*time.Minute); errKey != nil {
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

	result, err := database.GetCache("products_backend")

	db, _ := database.GetDatabaseConnection()

	if err != nil {
		db.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		if errKey := database.SetCache("products_backend", bytes, 30*time.Minute); errKey != nil {
			panic(errKey)
		}
	}

	if result != "" {
		json.Unmarshal([]byte(result), &products)
	}

	var searchedProducts []models.Product

	searchedProducts = products

	if search := c.Query("search"); search != "" {
		lowerCase := strings.ToLower(search)

		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), lowerCase) || strings.Contains(strings.ToLower(product.Description), lowerCase) {
				searchedProducts = append(searchedProducts, product)
			}
		}
	}

	if sortParam := c.Query("sort"); sortParam != "" {
		sortLower := strings.ToLower(sortParam)

		if sortLower == "asc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[i].Price < searchedProducts[j].Price
			})
		}

		if sortLower == "desc" {
			sort.Slice(searchedProducts, func(i, j int) bool {
				return searchedProducts[j].Price < searchedProducts[i].Price
			})
		}
	}

	var total = len(searchedProducts)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("perPage", "10"))

	var data []models.Product = []models.Product{}

	if total <= page*perPage && total >= (page-1)*perPage {
		data = searchedProducts[(page-1)*perPage : total]
	}

	if total >= page*perPage {
		data = searchedProducts[(page-1)*perPage : page*perPage]
	}

	return c.JSON(fiber.Map{
		"data":      data,
		"total":     total,
		"page":      page,
		"perPage":   perPage,
		"last_page": total/perPage + 1,
	})
}
