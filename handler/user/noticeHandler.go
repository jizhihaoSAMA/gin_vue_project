package user

import (
	"gin_vue_project/common"
	"gin_vue_project/dto"
	"gin_vue_project/model"
	"gin_vue_project/response"
	"gin_vue_project/service/userService/userNotice"
	"github.com/gin-gonic/gin"
	"sort"
	"time"
)

type Notice interface {
	GetTime() time.Time
}

func GetRecentNotices(ctx *gin.Context) {
	// 从所有通知中拉取最近10个通知。

	// 此处要修改为从请求头获取

	//userIDInString:= ctx.PostForm("user_id")
	userIDInString, _ := ctx.Get("isUser")
	userID, ok := userIDInString.(uint)
	if !ok {
		response.Fail(ctx, nil, "参数错误")
		return
	}

	recentVote := userNotice.GetRecentVotes(userID)
	recentComment := userNotice.GetRecentComment(userID)
	_ = userNotice.GetRecentFollowers(userID)

	db := common.InitMySQL()
	defer db.Close()

	var votesDto []dto.VoteDto

	// 获取comment，vote之后需要对结果进行处理
	for _, v := range recentVote {
		var voteDto dto.VoteDto
		var targetUser model.User
		var targetComment model.Comment

		voteDto.FormUserID = v.FromUserID
		voteDto.TargetCommentID = v.TargetCommentID
		voteDto.CreateAt = v.CreatedAt

		db.Model(&model.User{}).Where("id = ?", v.FromUserID).Scan(&targetUser)
		voteDto.Username = targetUser.Username
		db.Model(&model.Comment{}).Where("id = ?", v.TargetCommentID).Scan(&targetComment)
		voteDto.NewsID = targetComment.NewsID

		votesDto = append(votesDto, voteDto)
	}

	var noticeList []Notice
	for _, c := range recentComment {
		noticeList = append(noticeList, c)
	}

	for _, v := range recentVote {
		noticeList = append(noticeList, v)
	}
	sort.Slice(noticeList, func(i, j int) bool {
		return noticeList[i].GetTime().Before(noticeList[j].GetTime())
	})

	response.Success(ctx, gin.H{
		"notice": noticeList,
	}, "")
}

func GetUnreadAmount(ctx *gin.Context) {
	// 获取未读取的数量
	db := common.InitMySQL()
	defer db.Close()

	var targetUser model.User
	db.Model(&model.User{}).Where("id = ?", ctx.Query("user_id")).Scan(&targetUser)

	response.Success(ctx, gin.H{
		"amount": targetUser.CountOfUnread,
	}, "")
}

func GetAllNotices(ctx *gin.Context) {
	// 用户点击查看所有通知来查看所有的通知。
}
