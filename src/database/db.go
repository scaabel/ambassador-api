package database

import (
	"ambassador/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func Connect() {
	var err error
	db, err := gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/ambassador"), &gorm.Config{})

	if err != nil {
		panic("Could not connect to database")
	}

	dbConn = db
}

func GetDatabaseConnection() (*gorm.DB, error) {
	sqlDB, err := dbConn.DB()

	if err != nil {
		return dbConn, err
	}

	if err := sqlDB.Ping(); err != nil {
		return dbConn, err
	}

	return dbConn, nil
}

func AutoMigrate() error {
	db, connErr := GetDatabaseConnection()

	if connErr != nil {
		return connErr
	}

	err := db.AutoMigrate(&models.User{})

	return err
}
