package comment

import (
	"gin_vue_project/response"
	"gin_vue_project/service/commentService"
	"github.com/gin-gonic/gin"
)

func GetCommentAmountHandler(ctx *gin.Context) {
	response.Success(ctx, gin.H{
		"comment_amount": commentService.GetCommentAmount(ctx.PostForm("news_id")),
	}, "")
}
