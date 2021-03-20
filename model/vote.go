package model

import "time"

/*
   数据库中的投票表

   TargetCommentID: 目标的评论 ID
   FromUserID: 来自 xx 用户
   Status: 1 赞
           0 未赞 / 未踩
           -1 踩
*/

type Vote struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	TargetCommentID uint      `form:"target_comment_id" binding:"required" json:"target_comment_id"`

	// 应该根据用户的token 获取
	FromUserID uint `form:"from_user_id" binding:"required" json:"from_user_id"`

	// 由于status字段可以是0，但由于初始化会导致字段默认为0，gin的效验器的exists字段不存在。为不冲突，使用指针。
	// 见：https://github.com/gin-gonic/gin/issues/491

	Status *int `form:"status" binding:"required" json:"status"`
}

func (v Vote) GetTime() time.Time {
	return v.CreatedAt
}
