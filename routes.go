package main

import (
	"gin_vue_project/handler"
	"gin_vue_project/handler/comment"
	"gin_vue_project/handler/news"
	"gin_vue_project/handler/rpc"
	"gin_vue_project/handler/user"
	"gin_vue_project/middleware"
	"gin_vue_project/service/userService/userSecurity"
	"github.com/gin-gonic/gin"
)

func BindRoutes(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	// User Service
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("register", user.RegisterHandler)
		authGroup.POST("login", user.LoginHandler)
		authGroup.GET("info", middleware.UserServiceAuthHandler(), user.GetInfoHandler)
		authGroup.GET("updateToken", middleware.UserServiceAuthHandler(), user.UpdateTokenHandler)
	}
	// User Info
	r.POST("/api/post/userIcon", middleware.UserServiceAuthHandler(), user.UploadIconHandler)
	r.POST("/api/post/updateInfo", middleware.UserServiceAuthHandler(), user.UpdateInfoHandler)

	// Test Service
	r.GET("/api/test", handler.Test)
	r.POST("/api/test", handler.TestWithPost)

	// News Service
	r.GET("/api/get/news", news.GetNewsHandler)

	// Vote Service
	r.POST("/api/post/voteOnComment", comment.VoteOnCommentHandler)

	// Comment Service
	r.GET("/api/get/comments", middleware.IsUser(), comment.GetCommentsHandler)
	r.POST("/api/post/comment", middleware.UserServiceAuthHandler(), comment.PostCommentHandler)

	r.POST("/api/post/commentAmount", comment.GetCommentAmountHandler)

	r.POST("/api/post/getPage", comment.GetPageOfCommentHandler)

	// Security Service
	r.POST("/api/post/getCaptcha", user.SendMessageHandler)

	// 修改手机号时的服务。
	r.POST("/api/post/authOfChangeTelephone", middleware.UserServiceAuthHandler(), userSecurity.AuthOfChangeTelephone)
	r.POST("/api/post/updateUserTelephone", middleware.UserServiceAuthHandler(), userSecurity.UpdateUserTelephone)

	// RPC
	r.POST("/api/post/translate", rpc.TranslateHandler)

	return r
}
