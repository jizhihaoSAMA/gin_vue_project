package middleware

import (
	"fmt"
	"gin_vue_project/common"
	"github.com/gin-gonic/gin"
)

func checkCaptchaValid() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		telephone := ctx.PostForm("telephone")
		captchaCode := ctx.PostForm("captcha")
		db := common.InitRedis()
		val, err := db.Get(telephone).Result()
		if err != nil {
			fmt.Println("Err:", err)
			ctx.Abort()
			return
		}
		if val == captchaCode {
			// 验证码正确，下一个
			ctx.Next()
		}
	}
}
