package main

import (
	"gin_vue_project/handler"
	"gin_vue_project/handler/userSecurity"
	"gin_vue_project/middleware"
	"github.com/gin-gonic/gin"
)

func BindRoutes(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	// User Service
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("register", handler.Register)
		authGroup.POST("login", handler.Login)
		authGroup.GET("info", middleware.UserServiceAuthHandler(), handler.GetInfo)
	}

	// Test Service
	r.GET("/test", handler.Test)
	r.POST("/test", handler.TestWithPost)

	// News Service
	r.GET("/api/get/news", handler.GetNews)

	// Comment Service
	r.GET("/api/get/comments", handler.GetComments)
	r.POST("/api/post/comment", middleware.UserServiceAuthHandler(), handler.PostComment)

	// Security Service
	r.POST("/api/post/getCaptcha", handler.SendMessageHandler)
	// 修改手机号时的服务。
	r.POST("/api/post/authOfChangeTelephone", middleware.UserServiceAuthHandler(), userSecurity.AuthOfChangeTelephone)
	r.POST("/api/post/updateUserTelephone", middleware.UserServiceAuthHandler(), userSecurity.UpdateUserTelephone)

	// RPC
	r.POST("/api/post/translate", handler.Translate)

	return r
}
