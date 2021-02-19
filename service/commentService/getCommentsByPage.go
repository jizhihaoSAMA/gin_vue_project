package commentService

import (
	"gin_vue_project/common"
	"gin_vue_project/dto"
	"gin_vue_project/model"
)

func GetCommentsByPage(page int, newsID string, senderID interface{}) []dto.CommentDto {
	eachPageAmount := 10
	currentPage := page
	db := common.InitMySQL()
	defer db.Close()
	var results []model.Comment

	db.Model(&model.Comment{}).Where("news_id = ?", newsID).Offset(eachPageAmount * (currentPage - 1)).Limit(eachPageAmount).Scan(&results)

	var resultsDto []dto.CommentDto
	// results为该新闻所对应的评论
	for _, result := range results {
		var resultDto dto.CommentDto

		resultDto.ID = result.ID
		resultDto.UserID = result.UserID
		resultDto.Comment = result.Comment
		resultDto.Floor = result.Floor
		resultDto.CreatedAt = result.CreatedAt
		resultDto.Upvote = result.Upvote
		resultDto.Downvote = result.Downvote

		resultDto.TargetCommentContent = ""
		resultDto.TargetCommentUserID = 0
		resultDto.TargetCommentUsername = ""
		resultDto.TargetCommentID = 0

		var queryUser model.User
		db.Model(&model.User{}).Where("id = ?", result.UserID).Scan(&queryUser)
		resultDto.Username = queryUser.Username

		// 查询目标评论
		if result.TargetCommentID != 0 {
			var queryComment model.Comment
			db.Model(&model.Comment{}).Where("id = ?", result.TargetCommentID).Scan(&queryComment)
			resultDto.TargetCommentContent = queryComment.Comment
			resultDto.TargetCommentUserID = queryComment.UserID
			var targetUser model.User
			db.Model(&model.User{}).Where("id = ?", resultDto.TargetCommentUserID).Scan(&targetUser)
			resultDto.TargetCommentUsername = targetUser.Username
			resultDto.TargetCommentID = result.TargetCommentID
		}

		if senderID == 0 { // 如果是游客，则所有的投票状态全为0
			resultDto.VoteStatus = 0
		} else { // 如果是用户，查询该用户在评论投票的status
			var vote model.Vote
			db.Model(&model.Vote{}).Select("status").Where("target_comment_id = ? and from_user_id = ?", result.ID, senderID).Scan(&vote)
			if vote.Status == nil { // 判断空指针，如果用户没有对该评论进行操作，则在Vote表中是找不到结果的，此时为空指针
				resultDto.VoteStatus = 0
			} else {
				resultDto.VoteStatus = *vote.Status
			}
		}
		resultsDto = append(resultsDto, resultDto)
	}

	return resultsDto
}
