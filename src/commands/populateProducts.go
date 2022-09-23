package main

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"math/rand"

	"github.com/bxcodec/faker/v3"
)

func main() {
	database.Connect()

	for i := 0; i < 30; i++ {
		product := models.Product{
			Title:       faker.Name(),
			Description: faker.Sentence(),
			Price:       float64(rand.Intn(90) + 10),
			Image:       faker.URL(),
		}

		db, _ := database.GetDatabaseConnection()

		db.Create(&product)
	}

}
