package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username  string `json:"username" form:"username" gorm:"type:varchar(20);not null"`
	Telephone string `json:"telephone" form:"telephone" gorm:"varchar(110);not null;unique"`
	Password  string `json:"password" form:"password" gorm:"size:255;not null"`
}
