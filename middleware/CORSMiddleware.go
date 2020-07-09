package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 设置跨域请求的来源
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		// 设置跨域请求的缓存时间
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		// 设置通过访问的方法, "GET,POST,DELETE,PUT"或所有方法
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		// 设置允许的请求头
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		// 设置是否认证信息
		ctx.Writer.Header().Set("Access-Control-Credentials", "true")

		// 如果是Options请求(预检请求)，直接返回200
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	}
}
