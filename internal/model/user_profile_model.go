package model

import (
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	UserProfileId int64          `gorm:"primary_key; AUTO_INCREMENT" json:"-"`
	Provinsi      string         `json:"provinsi" form:"provinsi"`
	Kota          string         `json:"kota" form:"kota"`
	Kecamatan     string         `json:"kecamatan" form:"kecamatan"`
	Kelurahan     string         `json:"kelurahan" form:"kelurahan"`
	UserId        int64          `json:"-" `
}
