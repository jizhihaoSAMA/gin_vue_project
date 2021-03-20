package user

import (
	"gin_vue_project/common"
	"gin_vue_project/model"
	"gin_vue_project/response"
	"gin_vue_project/service/userService"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

func GetFollowUserHandler(ctx *gin.Context) {
	// 仍然是尝试分页。这个不急吧· - ·
}

func FollowUserHandler(ctx *gin.Context) {
	// 需要 来源的用户ID， 目标用户ID，以及对用户进行效验

	// 来源的ID在请求头里面，同时进行效验
	// 目标ID在POST方法中

	targetUserID := ctx.PostForm("follow_user_id")

	if targetUserID == "" { // 传递值为空
		response.Fail(ctx, nil, "请正确选择目标用户ID")
		return
	}

	var newFollowRow model.FollowRow

	tmpID, _ := ctx.Get("is_user")
	userID := tmpID.(uint)
	if userID == 0 {
		response.Fail(ctx, nil, "请登录")
	} else {
		newFollowRow.FromUserID = userID

		// 转为uint类型
		targetUserIDuint, _ := strconv.Atoi(targetUserID)
		newFollowRow.FollowUserID = uint(targetUserIDuint)
		db := common.InitMySQL()
		defer db.Close()

		// 创建记录
		db.Create(&newFollowRow)
		// 对用户的值进行操作
		db.Model(&model.User{}).Where("id = ?", targetUserID).Update("followers", gorm.Expr("followers + 1"))
		db.Model(&model.User{}).Where("id = ?", userID).Update("followings", gorm.Expr("followings + 1"))
	}

	response.Success(ctx, nil, "操作成功")
}

func UnfollowUserHandler(ctx *gin.Context) {
	// 需要 来源的用户ID， 目标用户ID，以及对用户进行效验

	// 来源的ID在请求头里面，同时进行效验
	// 目标ID在POST方法中

	tmpID, _ := ctx.Get("is_user")
	userID := tmpID.(uint)

	if userID == 0 {
		response.Fail(ctx, nil, "请登录")
		return
	}

	targetUserIDInForm := ctx.PostForm("unfollow_user_id")
	targetUserIDInInt, err := strconv.Atoi(targetUserIDInForm)
	if err != nil {
		response.ServerError(ctx, nil, "服务器出错：Err 001")
		return
	}
	targetUserID := uint(targetUserIDInInt)

	if userService.CheckUserExistsWithID(targetUserID) { // 传递值为空
		response.Fail(ctx, nil, "目标用户不存在")
		return
	}

	if targetUserID == userID {
		response.Fail(ctx, nil, "不能关注自己")
		return
	} else {
		// 删除follow表中的记录，并减去user表中follower的值，以及请求用户的值
		db := common.InitMySQL()
		defer db.Close()

		// 删除记录
		db.Unscoped().Where("from_user_id = ? and follow_user_id = ?", userID, targetUserID).Delete(&model.FollowRow{})

		// 对用户的值进行操作
		db.Model(&model.User{}).Where("id = ?", targetUserID).Update("followers", gorm.Expr("followers - 1"))
		db.Model(&model.User{}).Where("id = ?", userID).Update("followings", gorm.Expr("followings - 1"))
	}
	response.Success(ctx, nil, "操作成功")
}
