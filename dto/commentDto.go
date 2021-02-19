package dto

import (
	"time"
)

type CommentDto struct {
	ID                    uint      `json:"id"`
	CreatedAt             time.Time `json:"created_at"`
	Comment               string    `json:"comment"`
	UserID                uint      `json:"user_id"`
	Username              string    `json:"username"`
	Floor                 int       `json:"floor"`
	VoteStatus            int       `json:"vote_status"`
	Upvote                int       `json:"upvote"`
	Downvote              int       `json:"downvote"`
	TargetCommentID       uint      `json:"target_comment_id"`
	TargetCommentContent  string    `json:"target_comment_content"`
	TargetCommentUserID   uint      `json:"target_comment_user_id"`
	TargetCommentUsername string    `json:"target_comment_username"`
}
