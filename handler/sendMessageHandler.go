package handler

import (
	"fmt"
	"gin_vue_project/common"
	"gin_vue_project/utils"
	"github.com/gin-gonic/gin"
	"time"
)

func SendMessageHandler(ctx *gin.Context) {
	telephone := ctx.PostForm("telephone")
	fmt.Println(telephone)
	// 如果电话为空
	if telephone == "" {
		ctx.JSON(200, gin.H{
			"code": 403,
			"msg":  "电话号码不能为空",
		})
	}
	captchaCode := utils.SendMessage(telephone)

	db := common.InitRedis()
	defer db.Close()
	// 判断该IP是否到达上限， 5次机会
	ctx.ClientIP()

	// 存入redis数据库
	// 1分钟过期
	err := db.Set(telephone, captchaCode, 1*time.Minute).Err()
	if err != nil {
		fmt.Printf("set score failed, err:%v\n", err)
		return
	}
}
