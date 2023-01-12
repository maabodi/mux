package model

import (
	"net/textproto"
	"time"

	"gorm.io/gorm"
)

type News struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	NewsId    int64          `gorm:"primary_key; AUTO_INCREMENT" json:"-"`
	Title     string         `json:"title" form:"title"`
	Content   string         `json:"content" form:"content"`
	Thumbnail string         `json:"thumbnail" form:"thumbnail"`
}
type FileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
	// contains filtered or unexported fields
}
