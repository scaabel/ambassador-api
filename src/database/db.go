package database

import (
	"ambassador/src/models"
	"fmt"

	"ambassador/src/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func GetDsn() string {

	return fmt.Sprintf(
		"%s@tcp(%s:%s)/%s",
		config.Config("DB_USER"),
		config.Config("DB_URL"),
		config.Config("DB_PORT"),
		config.Config("DB_NAME"),
	)
}

func Connect() {

	var err error

	dsn := GetDsn()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

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

	err := db.AutoMigrate(&models.User{}, &models.Product{})

	return err
}
