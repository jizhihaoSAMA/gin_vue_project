package userNotice

import (
	"gin_vue_project/common"
	"gin_vue_project/model"
)

func GetRecentComment(userID uint) []model.Comment {
	db := common.InitMySQL()
	defer db.Close()

	var result []model.Comment
	db.Model(&model.Comment{}).Joins("inner join comments as father on father.id = comments.target_comment_id").Where("father.user_id = ?", userID).Scan(&result)

	return result
}
