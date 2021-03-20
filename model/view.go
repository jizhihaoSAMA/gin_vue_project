package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type View struct {
	gorm.Model
	NewsID               string `json:"news_id"`
	UserID               uint   `json:"user_id"`
	ViewAmountUpdateTime time.Time
}
