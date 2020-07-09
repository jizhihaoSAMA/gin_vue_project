package main

import (
	"gin_vue_project/handler"
	"gin_vue_project/middleware"
	"github.com/gin-gonic/gin"
)

func BindRoutes(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.POST("/api/auth/register", handler.Register)
	r.POST("/api/auth/login", handler.Login)
	r.GET("/api/auth/info", middleware.AuthHandler(), handler.GetInfo)
	return r
}
