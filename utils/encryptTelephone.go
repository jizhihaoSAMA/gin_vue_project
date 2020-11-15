package utils

import (
	"gin_vue_project/model"
	"log"
)

func Trans(userInterface interface{}) model.User {
	p, ok := (userInterface).(model.User)
	if ok {
		log.Println("成功")
		return p
	} else {
		log.Println("失败")
		return model.User{}
	}
}
