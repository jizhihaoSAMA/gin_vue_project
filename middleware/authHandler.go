package middleware

import (
	"gin_vue_project/common"
	"gin_vue_project/dto"
	"gin_vue_project/model"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		// 验证其是否合法
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(401, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)

		if err != nil || !token.Valid {
			ctx.JSON(401, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		// 验证通过后获取claim中的userId

		userId := claims.UserId
		db := common.InitMySQL()
		var user model.User
		db.First(&user, userId)

		// 用户查找 但不存在
		if userId == 0 {
			ctx.JSON(401, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		// 将用户的存在写入上下文
		ctx.Set("user", dto.ToUserDto(user))
		ctx.Next()

	}
}
