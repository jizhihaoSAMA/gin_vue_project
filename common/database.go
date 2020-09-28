package common

import (
	"context"
	"fmt"
	"gin_vue_project/model"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func InitMySQL() *gorm.DB {
	driverName := viper.GetString("dataSource.MySQL.driverName")
	host := viper.GetString("dataSource.MySQL.host")
	port := viper.GetString("dataSource.MySQL.port")
	database := viper.GetString("dataSource.MySQL.database")
	username := viper.GetString("dataSource.MySQL.username")
	password := viper.GetString("dataSource.MySQL.password")
	charset := viper.GetString("dataSource.MySQL.charset")
	loc := viper.GetString("dataSource.MySQL.loc")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		loc,
	)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("链接失败，错误:" + err.Error())
	}
	db.AutoMigrate(&model.User{}, &model.Comment{})
	return db
}

func InitMongoDB() (*mongo.Database, context.CancelFunc) {
	uri := fmt.Sprintf("%s://%s:%s",
		viper.GetString("dataSource.MongoDB.driverName"),
		viper.GetString("dataSource.MongoDB.host"),
		viper.GetString("dataSource.MongoDB.port"),
	)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	// 创建一个有超时限制的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(viper.GetString("dataSource.MongoDB.database")), cancel
}

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "root", // no password set
		DB:       0,      // use default DB
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		return nil
	}
	return rdb
}
