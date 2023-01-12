package model

import (
	"time"

	"gorm.io/gorm"
)

type Avatar struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	AvatarId  int64          `gorm:"primaryKey;AUTO_INCREMENT" json:"-"`
	AvatarUrl string         `json:"avatar_url" form:"avatar_url"`
	UserId    int64          `json:"-"`
}
