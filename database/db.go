package database

import (
	"content_portal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	HOST_POSTGRES = "34.128.93.96"
	PORT_POSTGRES = 5432
	DB_POSTGRES   = "halocat_db"
	USER_POSTGRES = "postgres"
	PASS_POSTGRES = "Noerhick_02"
)

var (
	db  *gorm.DB
	err error
)

func ConnectDB() *gorm.DB {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST_POSTGRES, PORT_POSTGRES, USER_POSTGRES, PASS_POSTGRES, DB_POSTGRES,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Debug().AutoMigrate(models.User{}, models.Company{}, models.Content{}, models.Point{})

	return db
}
