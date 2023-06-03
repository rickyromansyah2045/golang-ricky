package database

import (
	"content_portal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	HOST_POSTGRES = "localhost"
	PORT_POSTGRES = 3306
	DB_POSTGRES   = "cms_contents"
	USER_POSTGRES = "root"
	PASS_POSTGRES = ""
)

var (
	db  *gorm.DB
	err error
)

func ConnectDB() *gorm.DB {

	dsn := "root:@tcp(127.0.0.1:3306)/cms_contents?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Debug().AutoMigrate(models.User{}, models.Company{}, models.Content{}, models.Point{})

	return db
}
