package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Content struct {
	Id         uint       `gorm:"primaryKey" json:"id"`
	Title      string     `gorm:"not null" json:"title" valid:"required~title is required"`
	IsiContent string     `json:"isi_content"`
	PhotoUrl   string     `gorm:"not null" json:"photo_url" valid:"required~photo_url is required"`
	UserId     uint       `json:"user_id"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	// Comment   []Comment  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	User *User
}

func (p *Content) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)
	if errCreate != nil {
		return errCreate
	}

	return
}

func (p *Content) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)
	if errCreate != nil {
		return errCreate
	}

	return
}
