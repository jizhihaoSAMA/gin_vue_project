package userNotice

import (
	"gin_vue_project/common"
	"gin_vue_project/model"
	"gin_vue_project/service/userService"
	"github.com/jinzhu/gorm"
)

func SetNoticeForUser(userID uint) bool {
	db := common.InitMySQL()
	defer db.Close()
	if userService.CheckUserExistsWithID(userID) {
		return false
	}
	db.Model(model.User{}).Where("id = ?", userID).Update("count_of_unread", gorm.Expr("count_of_unread + 1"))
	return true
}

func SetNoticeForUserByComment(comment model.Comment) bool {
	db := common.InitMySQL()
	defer db.Close()

	var targetComment model.Comment
	// 查询到目标用户ID
	db.Model(model.Comment{}).Where("ID = ?", comment.TargetCommentID).Scan(&targetComment)
	if targetComment.UserID == comment.UserID { // 如果用户自己评论自己，则不会增加通知
		return true
	} else { // 否则尝试给用户设置通知
		return SetNoticeForUser(comment.UserID)
	}

}
