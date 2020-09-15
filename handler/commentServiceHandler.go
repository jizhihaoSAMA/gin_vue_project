package handler

import (
	"gin_vue_project/common"
	"gin_vue_project/dto"
	"gin_vue_project/model"
	"gin_vue_project/utils"
	"github.com/gin-gonic/gin"
	"log"
)

func PostComment(ctx *gin.Context) {

	var comment model.Comment
	err := ctx.Bind(&comment)

	userDtoInterface, _ := ctx.Get("user")
	userDto := utils.InterfaceToUserDto(userDtoInterface)

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
	if len(comment.Comment) > 300 || len(comment.Comment) < 5 {
		ctx.JSON(400, gin.H{
			"code":  400,
			"error": "Too long or too short.",
		})
		return
	}

	db := common.InitMySQL()
	defer db.Close()
	_ = db.Create(&comment)
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
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": result,
	})
}
