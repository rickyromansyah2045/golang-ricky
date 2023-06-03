package models

import (
	"content_portal/helpers"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Company struct {
	Id        uint       `gorm:"primaryKey" json:"id"`
	Username  string     `gorm:"not null;uniqueIndex" json:"username" valid:"required~username is required"`
	Email     string     `gorm:"not null;uniqueIndex" json:"email" valid:"required~email is required"`
	Password  string     `gorm:"not null" json:"password" valid:"required~password is required"`
	Role      bool       `gorm:"not null" json:"role" valid:"required~role is required"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (u *Company) BeforeCreate(tx *gorm.DB) (err error) {
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

func (u *Company) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)
	if errCreate != nil {
		return errCreate
	}

	return
}
