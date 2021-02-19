package comment

import (
	"gin_vue_project/response"
	"gin_vue_project/service/commentService"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetCommentsHandler(ctx *gin.Context) {
	page := ctx.Query("page")
	if page == "" { // 未传递时默认为第1页
		page = "1"
	}
	currentPage, err := strconv.Atoi(page)
	if err != nil { // 传递为空或者为字母
		response.Fail(ctx, nil, "page参数有误")
		return
	}

	senderID, ok := ctx.Get("isUser")

	if ok == false {
		response.ServerError(ctx, nil, "传入的参数有误，isUser参数未在上下文。")
		return
	}
	response.Success(ctx, gin.H{
		"comment_info": gin.H{
			"comments":     commentService.GetCommentsByPage(currentPage, ctx.Query("news_id"), senderID),
			"current_page": page,
		},
	}, "")

}
