package dto

import "time"

type VoteDto struct {
	FormUserID      uint      `json:"from_user_id"`
	Username        string    `json:"username"`
	TargetCommentID uint      `json:"target_comment_id"`
	NewsID          string    `json:"news_id"`
	CreateAt        time.Time `json:"create_at"`
}

func (v VoteDto) GetTime() time.Time {
	return v.CreateAt
}
