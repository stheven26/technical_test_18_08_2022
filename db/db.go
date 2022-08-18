package db

import (
	"github.com/stheven26/config"
	"github.com/stheven26/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func SetupDB() {
	config := config.ConfigDB()

	dsn := config.DB_USERNAME + ":" + config.DB_PASSWORD + "@(" + config.DB_HOST + ")/" + config.DB_NAME + "?charset=utf8&parseTime=True&loc=Local"

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Blog{})
}

func GetConnectionDB() *gorm.DB {
	return db
}
