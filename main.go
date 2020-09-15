package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	InitConfig()
	r := gin.Default()
	r = BindRoutes(r)

	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}

	panic(r.Run())
}

func InitConfig() {
	workDir, _ := os.Getwd()
	log.SetFlags(log.Ldate | log.Lshortfile)
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
