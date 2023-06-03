package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Point struct {
	Id        uint       `gorm:"primaryKey" json:"id"`
	UserId    uint       `json:"user_id"`
	ContentId uint       `json:"content_id"`
	Point     uint       `gorm:"not null" json:"point" valid:"required~point is required"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	User    *User
	Content *Content
}

func (c *Point) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)
	if errCreate != nil {
		return errCreate
	}
	return

}

func (c *Point) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(c)
	if errCreate != nil {
		return errCreate
	}

	return
}
