package userService

import (
	"gin_vue_project/common"
	"gin_vue_project/model"
	"log"
)

func CheckUserExistsWithID(userID uint) bool {
	db := common.InitMySQL()
	defer db.Close()
	var user model.User
	db.Model(&model.User{}).Where("id = ?", userID).Scan(&user)
	log.Printf("%+v", user)
	if userID == 0 {
		return false
	} else {
		return true
	}
}
