package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	RoleId    int64          `gorm:"primaryKey" json:"-"`
	RoleName  string         `json:"name" form:"name"`
}
