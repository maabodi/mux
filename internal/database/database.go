package database

import (
	"muxapp/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost)/gomux?charset=utf8&parseTime=true"))

	if err != nil {
		panic(err)
	}

	db.Create(&model.Role{
		RoleId:   1,
		RoleName: "User",
	})
	db.Create(&model.Role{
		RoleId:   2,
		RoleName: "Operator",
	})
	db.Create(&model.Role{
		RoleId:   3,
		RoleName: "Admin",
	})

	db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.News{},
		&model.UserProfile{},
		&model.Avatar{},
		&model.Activity{},
		&model.WProvinsi{},
		&model.WKota{},
		&model.WKecamatan{},
		&model.WKelurahan{},
	)

	DB = db
}
