package userSecurity

import (
	"github.com/gin-gonic/gin"
)

func AuthOfChangeTelephone(ctx *gin.Context) {
	//userID := ctx.PostForm("userID")
	//captcha := ctx.PostForm("captcha")
	//formTelephone := ctx.PostForm("formTelephone")
	//
	//// 验证码检测
	//rdb := common.InitRedis()
	//defer rdb.Close()
	//
	//captchaInRedis, _ := rdb.Get(formTelephone+"_sec").Result()
	//if captchaInRedis == captcha { // 相同代表正确
	//    // 增加session。
	//    ctx.JSON(200, gin.H{
	//        "code": "200",
	//        "msg": "操作成功",
	//    })
	//} else {
	//    ctx.JSON(403, gin.H{
	//        "code": "403",
	//        "msg":  "验证码不正确",
	//    })
	//    return
	//}

}
