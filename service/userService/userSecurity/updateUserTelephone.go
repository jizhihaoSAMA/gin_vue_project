package userSecurity

import (
	"gin_vue_project/common"
	"gin_vue_project/model"
	"gin_vue_project/response"
	"github.com/gin-gonic/gin"
	"log"
)

func UpdateUserTelephone(ctx *gin.Context) {
	userID := ctx.PostForm("userID")
	captcha := ctx.PostForm("captcha")
	formTelephone := ctx.PostForm("formTelephone")

	// 先检查手机号是否正确
	mysqlDB := common.InitMySQL()
	defer mysqlDB.Close()
	var user model.User
	mysqlDB.Where("id = ?", userID).First(&user)
	if user.Telephone == formTelephone { // 表单的电话正确，方可进入下一步
		redisDB := common.InitRedis()
		defer redisDB.Close()

		// 检查验证码是否正确
		captchaInRedis, _ := redisDB.Get(formTelephone + "_sec").Result()
		if captchaInRedis == captcha { // 验证码相同
			// 原手机号 + 验证码都正确，修改手机号
			response.Success(ctx, nil, "操作成功")
			// 成功删除次数以及号码
			log.Println(redisDB.Del(formTelephone+"_sec"), redisDB.Del(formTelephone+"_sec_times"))

		} else { // 不正确。
			remainTimes, _ := redisDB.Decr(formTelephone + "_sec_times").Result()

			log.Println(remainTimes, remainTimes)

			if remainTimes == -1 {
				response.Fail(ctx, nil, "验证码失效")
				// 删除该验证码以及剩余次数
				log.Println(redisDB.Del(formTelephone+"_sec"), redisDB.Del(formTelephone+"_sec_times"))
			} else if remainTimes >= 0 {
				response.Fail(ctx, nil, "验证码错误，请检查输入是否正确。")
			} else {
				response.Response(ctx, 500, 500, nil, "内部错误")
			}

		}
	} else {
		response.Fail(ctx, nil, "电话号码错误")
	}

}
