package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	UserID      int64          `gorm:"primary_key; AUTO_INCREMENT" json:"-"`
	Name        string         `json:"name" form:"name"`
	Email       string         `json:"email" form:"email"`
	Password    string         `json:"password" form:"password"`
	RoleId      int64          `json:"role"`
	Role        Role           `json:"-"`
	UserProfile UserProfile    `json:"user_profile"`
	Avatar      Avatar         `gorm:"constraint:OnUpdate:CASCADE,ONDELETE:SET NULL;" json:"avatar"`
}

type ResUser struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar Avatar `json:"avatar"`
}
