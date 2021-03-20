package handler

import (
	"gin_vue_project/common"
	"gin_vue_project/model"
	"gin_vue_project/response"
	"github.com/gin-gonic/gin"
)

func Test(ctx *gin.Context) {
	// 获取未读消息， 1，点赞 2，关注 3，评论
	userID := ctx.Query("user_id")

	db := common.InitMySQL()
	defer db.Close()
	// 先获取点赞的前10个消息。

	type result struct {
		FromUserID      uint
		Status          int
		TargetCommentID uint
	}
	var resultList []result

	//db.LogMode(true)

	// 该段注释为不给出
	//db.Model(&model.Vote{}).Select("votes.from_user_id, votes.status, votes.target_comment_id").Joins("left join comments on comments.id = votes.target_comment_id").Where("comments.user_id = ? and comments.user_id != votes.from_user_id", userID).Limit(10).Scan(&resultList)
	db.Model(&model.Vote{}).Select("votes.from_user_id, votes.status, votes.target_comment_id").Joins("left join comments on comments.id = votes.target_comment_id").Where("comments.user_id = ?", userID).Limit(10).Scan(&resultList)

	response.Success(ctx, gin.H{
		"data": resultList,
	}, "")
}

func TestWithPost(ctx *gin.Context) {
	commentID := ctx.PostForm("comment_id")
	db := common.InitMySQL()
	defer db.Close()

	var targetComment model.Comment
	db.Model(&model.Comment{}).Where("id = ?", commentID).Scan(&targetComment)

	var pos int
	// 该评论在页码的位置
	db.Model(&model.Comment{}).Where("news_id = ? and id <= ?", targetComment.NewsID, commentID).Count(&pos)

	pageNumber := (pos / 10) + 1

	response.Success(ctx, gin.H{
		"comment_info": gin.H{
			"comments":     nil,
			"current_page": pageNumber,
		},
	}, "")

}
