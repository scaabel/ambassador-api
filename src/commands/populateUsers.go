package main

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/bxcodec/faker/v3"
)

func main() {
	database.Connect()

	for i := 0; i < 30; i++ {
		ambassador := models.User{
			Name:         faker.Name(),
			Email:        faker.Email(),
			IsAmbassador: true,
		}

		ambassador.SetPassword("1234")

		db, _ := database.GetDatabaseConnection()

		db.Create(&ambassador)
	}

}
