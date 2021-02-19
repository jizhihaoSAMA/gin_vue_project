package comment

import (
	"context"
	"gin_vue_project/common"
	"gin_vue_project/dto"
	"gin_vue_project/model"
	"gin_vue_project/response"
	"gin_vue_project/service/userService/userNotice"
	"gin_vue_project/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func PostCommentHandler(ctx *gin.Context) {

	var comment model.Comment
	var normalNews model.NormalNews
	err := ctx.ShouldBind(&comment)

	if err != nil {
		response.Fail(ctx, nil, "提交数据错误")
		return
	}

	userDtoInterface, _ := ctx.Get("user")
	userDto := utils.InterfaceToUserDto(userDtoInterface)
	// 判断用户是否为空，即用户是否合法
	if userDto == (dto.UserDto{}) {
		ctx.JSON(401, gin.H{
			"code": 401,
			"msg":  "Auth Failed",
		})
		return
	}

	comment.UserID = userDto.ID

	log.Println(comment.TargetCommentID)
	// 判断评论长度
	if len(comment.Comment) >= 50 || len(comment.Comment) < 5 {
		ctx.JSON(400, gin.H{
			"code":  400,
			"error": "Too long or too short.",
		})
		return
	}

	// 条件符合，插入数据，先获取当前新闻的comment数量，使其数量+1，+1后的结果为当前评论的楼层
	db, cancel := common.InitMongoDB()
	defer cancel()

	id, err := primitive.ObjectIDFromHex(comment.NewsID)

	filter := bson.D{{"_id", id}}

	update := bson.D{{
		"$inc", bson.D{{
			"comment_amount", 1,
		}},
	}}

	after := options.After

	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	err = db.Collection("test").FindOneAndUpdate(context.TODO(), filter, update, &opt).Decode(&normalNews)

	if err != nil {
		log.Fatal(err)
	}

	comment.Floor = normalNews.CommentAmount

	mysqlDB := common.InitMySQL()
	defer mysqlDB.Close()
	_ = mysqlDB.Create(&comment)

	if !userNotice.SetNoticeForUserByComment(comment) {
		response.Fail(ctx, nil, "用户不存在，无法添加通知")
		return
	}

	response.Success(ctx, nil, "操作成功")
}
