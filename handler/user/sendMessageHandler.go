package user

import (
	"fmt"
	"gin_vue_project/common"
	"gin_vue_project/model"
	"gin_vue_project/response"
	"gin_vue_project/utils"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func SendMessageHandler(ctx *gin.Context) {
	telephone := ctx.PostForm("telephone")
	log.Println(time.Now(), telephone)
	// 如果电话为空
	if telephone == "" {
		ctx.JSON(200, gin.H{
			"code": 403,
			"msg":  "电话号码不能为空",
		})
	}

	redisDB := common.InitRedis()
	defer redisDB.Close()
	mysqlDB := common.InitMySQL()
	defer mysqlDB.Close()
	// 判断该IP是否到达上限， 5次机会
	//ctx.ClientIP()

	// 存入redis数据库
	// 1分钟过期
	if ctx.Query("for") == "register" {
		// 判断用户是否存在
		if isTelephoneExist(mysqlDB, telephone) {
			response.Response(ctx, 400, 400, nil, "用户已存在，请更换手机号注册")
			return
		}
		captchaCode := utils.SendMessage(telephone)
		err := redisDB.Set(telephone+"_reg", captchaCode, 5*time.Minute).Err()
		if err != nil {
			fmt.Printf("set \"telephone for register\" failed, err:%v\n", err)
			return
		}
	} else if ctx.Query("for") == "security" {
		// 检查电话号码是否是当前用户所绑定的号码
		userID := ctx.PostForm("userID")
		if userID == "" {
			response.Response(ctx, 405, 405, nil, "提交信息有误")
			return
		}
		mysqlDB := common.InitMySQL()
		defer mysqlDB.Close()
		var user model.User
		mysqlDB.Where("id = ?", userID).First(&user)
		if user.ID == 0 || user.Telephone != ctx.PostForm("telephone") { // 查询不到结果, 或结果不符合
			response.Response(ctx, 403, 403, nil, "输入有误，请检查手机号是否正确")
			return
		}
		// 数据成功，发送手机号
		captchaCode := utils.SendMessage(telephone)
		err := redisDB.Set(telephone+"_sec", captchaCode, 5*time.Minute).Err()
		if err != nil {
			fmt.Printf("set \"telephone for security\" failed, err:%v\n", err)
			return
		}
		// 5次试错机会。
		err = redisDB.Set(telephone+"_sec_times", "5", 5*time.Minute).Err()
		if err != nil {
			fmt.Printf("set \"telephone times for security\" failed, err:%v\n", err)
			return
		}
		ctx.JSON(200, gin.H{
			"code": "200",
			"msg":  "发送成功",
		})
		return
	} else {
		fmt.Println("Another something")
	}
}
