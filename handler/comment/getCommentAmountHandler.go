package comment

import (
	"gin_vue_project/response"
	"gin_vue_project/service/commentService"
	"github.com/gin-gonic/gin"
	"log"
)

func GetCommentAmountHandler(ctx *gin.Context) {
	log.Println(ctx.PostForm("news_id"))
	response.Success(ctx, gin.H{
		"comment_amount": commentService.GetCommentAmount(ctx.PostForm("news_id")),
	}, "")
}
