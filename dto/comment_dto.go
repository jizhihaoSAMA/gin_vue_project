package dto

import (
	"time"
)

type Comment struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Comment   string    `json:"comment"`
	UserID    uint      `json:"user_id"`
	Floor     uint      `json:"floor"`
}
