package newsService

import (
	"gin_vue_project/common"
	"github.com/go-redis/redis"
	"log"
)

func AddHotValueByNewsID(newsID string, value float64) bool {
	redisDB := common.InitRedis()
	defer redisDB.Close()
	// 获取热度

	newsHotValue, _ := redisDB.ZScore("hot_news", newsID).Result()
	err := redisDB.ZAdd("hot_news", redis.Z{Score: newsHotValue + value, Member: newsID}).Err()

	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
