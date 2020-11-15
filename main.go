package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"os"
)

func main() {
	InitConfig()
	gin.DisableConsoleColor()
	r := gin.Default()
	r = BindRoutes(r)

	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "root", []byte("secret"))
	r.Use(sessions.Sessions("user", store))

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
