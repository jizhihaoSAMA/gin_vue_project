package main

import (
	"gin_vue_project/common"
	"gin_vue_project/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mattn/go-colorable"
	"github.com/spf13/viper"
	"os"
)

func main() {
	InitConfig()
	gin.ForceConsoleColor()
	gin.DefaultWriter = colorable.NewColorableStdout()
	gin.DefaultErrorWriter = colorable.NewColorableStderr()
	r := gin.Default()
	r = BindRoutes(r)

	//r.StaticFS("/user/userInfo", http.Dir("static"))
	r.Static("api/userIcon", "./static/userInfo/userIcon")
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "root", []byte("secret"))
	r.Use(sessions.Sessions("user", store))

	// 启动时自动更新一下
	db := common.InitMySQL()
	db.AutoMigrate(
		&model.User{},
		&model.Comment{},
		&model.Vote{},
		&model.FollowRow{},
		&model.View{},
	)

	db.LogMode(true)

	port := viper.GetString("port.server")
	if port != "" {
		panic(r.Run(":" + port))
	}

	panic(r.Run())
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
