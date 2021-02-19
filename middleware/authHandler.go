package middleware

import (
	"gin_vue_project/common"
	"gin_vue_project/dto"
	"gin_vue_project/model"
	"github.com/gin-gonic/gin"
	"strings"
)

// 对用户的登录进行验证的中间件，直接拦截
func UserServiceAuthHandler() gin.HandlerFunc {
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

// 验证请求方是游客还是用户
/*
	如果是游客，则用户ID为0，否则值是用户ID。
*/
func IsUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		// 验证其是否合法
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.Set("isUser", 0)
			ctx.Next()
			return
		}

		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)

		if err != nil || !token.Valid {
			ctx.Set("isUser", 0)
			ctx.Next()
			return
		}

		// 验证通过后获取claim中的userId

		userId := claims.UserId
		db := common.InitMySQL()
		var user model.User
		db.First(&user, userId)

		// 将用户的ID写入上下文
		ctx.Set("isUser", dto.ToUserDto(user).ID)
		ctx.Next()
	}
}

func UserSessionValid() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
