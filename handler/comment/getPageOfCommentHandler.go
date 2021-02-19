package comment

import (
	"gin_vue_project/common"
	"gin_vue_project/model"
	"gin_vue_project/response"
	"github.com/gin-gonic/gin"
)

func GetPageOfCommentHandler(ctx *gin.Context) {
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
