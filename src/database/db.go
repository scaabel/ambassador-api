package database

import (
	"ambassador/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() {
	var err error
	db, err := gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/ambassador"), &gorm.Config{})

	if err != nil {
		panic("Could not connect to database")
	}

	db.AutoMigrate(models.User{})
}
