package models

import (
	"content_portal/helpers"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	Id        uint       `gorm:"primaryKey" json:"id"`
	Username  string     `gorm:"not null;uniqueIndex" json:"username" valid:"required~username is required"`
	Email     string     `gorm:"not null;uniqueIndex" json:"email" valid:"required~email is required"`
	Password  string     `gorm:"not null" json:"password" valid:"required~password is required"`
	Age       int        `gorm:"not null" json:"age" valid:"required~age is required"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Content   []Content  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// Point     []Point    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Point *Point
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)
	if errCreate != nil {
		return errCreate
	}

	// validasi Password
	hash, err := helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = hash

	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)
	if errCreate != nil {
		return errCreate
	}

	return
}
