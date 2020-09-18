package handler

import (
	"context"
	"encoding/json"
	"gin_vue_project/common"
	"gin_vue_project/dto"
	"gin_vue_project/model"
	"gin_vue_project/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func PostComment(ctx *gin.Context) {

	var comment model.Comment
	var normalNews model.NormalNews
	err := ctx.Bind(&comment)

	userDtoInterface, _ := ctx.Get("user")
	userDto := utils.InterfaceToUserDto(userDtoInterface)
	// 判断用户是否为空
	if userDto == (dto.UserDto{}) {
		ctx.JSON(401, gin.H{
			"code": 401,
			"msg":  "Auth Failed",
		})
	}

	comment.UserID = userDto.ID
	if err != nil {
		log.Fatal(err)
	}
	// 判断评论长度
	if len(comment.Comment) > 300 || len(comment.Comment) < 5 {
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
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "Operate successfully",
	})
}

func GetComments(ctx *gin.Context) {
	db := common.InitMySQL()
	defer db.Close()
	var result []model.Comment
	db.Where("news_id = ? AND is_deleted = false", ctx.Query("news_id")).Find(&result)

	b, _ := json.Marshal(&result)
	var tmp []gin.H
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		log.Fatal(err)
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": tmp,
	})
}
