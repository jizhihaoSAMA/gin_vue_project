package model

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey;auto_increment" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	IsDeleted bool      `gorm:"bool;default:false"`
	DeletedAt time.Time `gorm:"default:null"`
	Comment   string    `form:"comment" gorm:"varchar(300)" json:"comment"`
	NewsID    string    `form:"news_id" json:"news_id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Floor     int       `json:"floor"`
}
