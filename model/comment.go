package model

import (
	"time"
)

/*
	gorm.Model

	Comment: 评论内容
	NewsID: 对应的news ID
	UserID： 发表的用户ID
	Username： Join查询产生的外键
	Floor：评论对应的楼层
	TargetUserID：评论的目标对象。默认为空，消息通知用
	TargetCommentID：评论的目标评论。默认为空，消息
	Upvote: 点赞数目。
 	Downvote: 踩 数目。
*/

type Comment struct {
	ID              uint      `json:"id"                     gorm:"primaryKey;auto_increment" `
	CreatedAt       time.Time `json:"created_at"             gorm:"autoCreateTime" `
	IsDeleted       bool      `json:"-"                      gorm:"bool;default:false" `
	DeletedAt       time.Time `json:"-"                      gorm:"default:null" sql:"index" `
	Comment         string    `json:"comment" form:"comment" gorm:"varchar(300)"`
	NewsID          string    `json:"news_id" form:"news_id"`
	UserID          uint      `json:"user_id"`
	Floor           int       `json:"floor"`
	Upvote          int       `json:"upvote"                 `
	Downvote        int       `json:"downvote"               `
	TargetCommentID uint      `json:"target_comment_id"      form:"target_comment_id"`
}
