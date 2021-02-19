package model

import (
	"github.com/jinzhu/gorm"
)

/*
   CommentID： 该回复所对应的Comment ID
   Reply： 评论内容
   ReplyID： 回复目标的ID

*/

type Reply struct {
	gorm.Model
	CommentID    int
	Reply        string `form:"reply" gorm:"varchar(300)" json:"reply"`
	ReplyID      int
	ReplyType    int
	FromUserID   int
	TargetUserID int
}
