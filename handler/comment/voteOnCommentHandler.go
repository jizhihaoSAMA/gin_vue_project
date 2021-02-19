package comment

import (
	"gin_vue_project/common"
	"gin_vue_project/model"
	"gin_vue_project/response"
	"gin_vue_project/service/userService"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
)

func VoteOnCommentHandler(ctx *gin.Context) {
	mapper := map[int]string{
		-1: "downvote",
		1:  "upvote",
	}

	var postVote model.Vote
	err := ctx.ShouldBind(&postVote)
	if err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据有误")
		return
	}
	if !userService.CheckUserExistsWithID(postVote.FromUserID) {
		response.Fail(ctx, nil, "From用户的ID不存在")
		return
	}
	var targetComment model.Comment
	//log.Printf("%+v\n", postVote)
	db := common.InitMySQL()
	defer db.Close()
	//db.Create(&vote)
	//
	//
	//db.Where("id = ?", vote.TargetCommentID).First(&comment)
	//log.Println(comment.Comment)
	if *postVote.Status == 1 {
		/*
		   1. 判断是否赞过，如果赞过，2不执行，直接返回已赞。

		   2. 找到目标评论，赞数+1

		   3. 点赞，删除之前的踩，则从vote表中找到与 target_comment_id 以及 from_user_id 一致的票
		       如果找不到用户的踩，继续。
		       如果找到之前的踩，则直接删除该vote记录，同时downvote数量减一。
		*/

		// step 1:
		var vote model.Vote
		// 查找用户是否赞过该评论
		db.Where("from_user_id = ? and status = 1 and target_comment_id = ?", postVote.FromUserID, postVote.TargetCommentID).First(&vote)
		if vote.ID != 0 { // 查到了upvote的结果代表已经赞过
			response.Fail(ctx, nil, "您已赞过。")
			return
		}

		// step 2:

		// 找到目标评论
		db.Where("id = ? and is_deleted = 0", postVote.TargetCommentID).First(&targetComment)
		if targetComment.ID == 0 { // 评论不存在或已经删除。
			response.Fail(ctx, nil, "评论不存在或已经删除。")
			return
		}
		// 值+1

		db.Model(&targetComment).Update("upvote", gorm.Expr("upvote + 1"))

		// step 3:
		vote = model.Vote{}
		// 查找用户是否踩过，删除之前的踩，且必须对应当前的comment
		db.Where("from_user_id = ? and status = -1 and target_comment_id = ?", postVote.FromUserID, postVote.TargetCommentID).First(&vote)
		if vote.ID != 0 {
			db.Unscoped().Delete(&model.Vote{}, vote.ID)
			db.Model(&targetComment).Update("downvote", gorm.Expr("downvote - 1"))
		}
		db.Save(&postVote)
		response.Success(ctx, nil, "点赞成功")

		// 通知用户，为用户增加一个通知，并将未读消息数量 + 1
		userService.SetNoticeForUserWithCommentID(targetComment.ID)

	} else if *postVote.Status == -1 {
		/*
		   1. 判断是否踩过，如果踩过，2不执行，直接返回已踩。

		   2. 找到目标评论，踩数+1

		   3. 点踩，删除之前的赞，则从vote表中找到与 target_comment_id 以及 from_user_id 一致的票
		       如果找不到用户的赞，继续。
		       如果找到之前的赞，则直接删除该vote记录，同时upvote数量减一。
		*/
		// step 1:
		var vote model.Vote
		// 查找用户是否踩过
		db.Where("from_user_id = ? and status = -1 and target_comment_id = ?", postVote.FromUserID, postVote.TargetCommentID).First(&vote)
		if vote.ID != 0 { // 查到了downvote的结果代表已经赞过
			response.Fail(ctx, nil, "您已踩过。")
			return
		}

		// step 2:

		// 找到目标评论
		db.Where("id = ? and is_deleted = 0", postVote.TargetCommentID).First(&targetComment)
		if targetComment.ID == 0 { // ID为初值，代表评论不存在或已经删除。
			response.Fail(ctx, nil, "评论不存在或已经删除。")
			return
		}
		// 值+1
		db.Model(&targetComment).Update("downvote", gorm.Expr("downvote + 1"))

		// step 3:
		vote = model.Vote{}
		// 查找用户是否赞过，删除之前的赞。
		db.Where("from_user_id = ? and status = 1 and target_comment_id = ?", postVote.FromUserID, postVote.TargetCommentID).First(&vote)
		if vote.ID != 0 {
			db.Unscoped().Delete(&model.Vote{}, vote.ID)
			db.Model(&targetComment).Update("upvote", gorm.Expr("upvote - 1"))
		}
		db.Save(&postVote)
		response.Success(ctx, nil, "点踩成功")

	} else if *postVote.Status == 0 { // 用户撤销投票
		// 先判断用户是否投过票
		var hadVoteBefore model.Vote

		// 查询语句有问题
		db.Model(&model.Vote{}).Where("from_user_id = ? and target_comment_id = ?", postVote.FromUserID, postVote.TargetCommentID).First(&hadVoteBefore)

		if hadVoteBefore.ID == 0 {
			response.Fail(ctx, nil, "用户并未投票，无需撤销")
			return
		}

		db.Where("id = ? and is_deleted = 0", postVote.TargetCommentID).First(&targetComment)
		if targetComment.ID == 0 { // ID为初值，代表评论不存在或已经删除。
			response.Fail(ctx, nil, "评论不存在或已经删除。")
			return
		}

		voteType := mapper[*hadVoteBefore.Status]
		db.Unscoped().Delete(hadVoteBefore)
		db.Model(&model.Comment{}).Where("id = ?", postVote.TargetCommentID).Update(voteType, gorm.Expr(voteType+" - 1"))

		response.Success(ctx, nil, "撤销成功")
	} else {
		response.Fail(ctx, nil, "status码错误。")
		return
	}
}
