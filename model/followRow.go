package model

import "github.com/jinzhu/gorm"

type FollowRow struct {
	gorm.Model
	FromUserID   uint `json:"from_user_id"`
	FollowUserID uint `json:"follow_user_id"`
}
