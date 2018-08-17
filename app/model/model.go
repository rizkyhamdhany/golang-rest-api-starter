package model

import (
	"github.com/jinzhu/gorm"
	"github.com/rizkyhamdhany/kelase-micro/app/module/admin"
	"github.com/rizkyhamdhany/kelase-micro/app/module/users"
)

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&admin.Admin{})
	db.AutoMigrate(&users.User{})
	return db
}
