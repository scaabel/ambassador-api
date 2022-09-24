package database

import (
	"ambassador/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func GetDsn() string {

	return "root@tcp(db:3306)/ambassador"
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

	err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Link{},
		&models.Order{},
		&models.OrderItem{},
	)

	return err
}
