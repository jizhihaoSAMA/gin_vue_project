package main

import (
	"gin_vue_project/handler"
	"gin_vue_project/middleware"
	"github.com/gin-gonic/gin"
)

func BindRoutes(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	// User Service
	r.POST("/api/auth/register", handler.Register)
	r.POST("/api/auth/login", handler.Login)
	r.GET("/api/auth/info", middleware.UserServiceAuthHandler(), handler.GetInfo)

	// Test Service
	r.GET("/test", handler.Test)
	r.POST("/test", handler.TestWithPost)

	// News Service
	r.GET("/api/get/news", handler.GetNews)

	// Comment Service
	r.GET("/api/get/comments", handler.GetComments)
	r.POST("/api/post/comment", middleware.UserServiceAuthHandler(), handler.PostComment)

	return r
}
