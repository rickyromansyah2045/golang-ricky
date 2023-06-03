package database

import (
	"content_portal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	HOST_POSTGRES = "34.101.130.27"
	PORT_POSTGRES = 3306
	DB_POSTGRES   = "halocat_db"
	USER_POSTGRES = "root"
	PASS_POSTGRES = "root123"
)

var (
	db  *gorm.DB
	err error
)

func ConnectDB() *gorm.DB {

	dsn := "root:root123@tcp(34.101.130.27:3306)/halocat_db?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Debug().AutoMigrate(models.User{}, models.Company{}, models.Content{}, models.Point{})

	return db
}
