package userNotice

import (
	"gin_vue_project/common"
	"gin_vue_project/model"
)

func GetRecentVotes(userID uint) []model.Vote {
	// 获取未读消息， 点赞

	db := common.InitMySQL()
	defer db.Close()
	// 先获取点赞的前10个消息。

	var resultList []model.Vote
	//db.LogMode(true)

	// 该段注释为不给出 自己给自己投票的通知
	//db.Model(&model.Vote{}).Select("votes.from_user_id, votes.status, votes.target_comment_id").Joins("left join comments on comments.id = votes.target_comment_id").Where("comments.user_id = ? and comments.user_id != votes.from_user_id", userID).Limit(10).Scan(&resultList)

	db.Model(&model.Vote{}).Joins("left join comments on comments.id = votes.target_comment_id").Where("comments.user_id = ?", userID).Limit(10).Scan(&resultList)
	return resultList
}
