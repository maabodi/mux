package model

import (
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	ActivityId   int64          `gorm:"primaryKey;AUTO_INCREMENT" json:"-"`
	ActivityName string         `json:"activity"`
	Place        string         `json:"place"`
	Time         time.Time      `json:"time"`
}
