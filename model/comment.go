package model

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey;auto_increment"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	IsDeleted bool      `gorm:"bool;default:false"`
	DeletedAt time.Time `gorm:"default:null"`
	Comment   string    `form:"comment" gorm:"varchar(300)"`
	NewsID    string    `form:"news_id"`
	UserID    uint
}
