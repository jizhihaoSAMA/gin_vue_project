package news

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_vue_project/common"
	"gin_vue_project/model"
	"gin_vue_project/response"
	"gin_vue_project/service/newsService"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func GetNewsHandler(ctx *gin.Context) {
	log.SetFlags(log.Ldate | log.Lshortfile)
	db, cancel := common.InitMongoDB()
	defer cancel()
	collection := db.Collection("test")

	id := ctx.Query("id")
	if id != "" { // 获取具体的新闻
		var newsContent model.NormalNews
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ctx.JSON(404, gin.H{
				"code": 404,
				"msg":  "Not found",
			})
		}
		filter := bson.M{"_id": objectID}
		var update bson.M
		var viewAdd int
		var targetView model.View

		_userID, _ := ctx.Get("isUser")
		userID, ok := _userID.(uint)
		if !ok {
			response.ServerError(ctx, nil, "用户ID出错")
			return
		}

		mysqlDB := common.InitMySQL()
		defer mysqlDB.Close()

		mysqlDB.Model(&model.View{}).Where("user_id = ? and news_id = ?", userID, id).Scan(&targetView)

		if targetView.ID == 0 { // 用户从未浏览过，创建访问记录并使浏览量 + 1
			targetView.NewsID = id
			targetView.UserID = userID
			targetView.ViewAmountUpdateTime = time.Now()
			mysqlDB.Create(&targetView)

			// 增加热度
			if ok := newsService.AddHotValueByNewsID(id, 1); !ok {
				response.ServerError(ctx, nil, "服务器错误")
			}

			viewAdd = 1
		} else { // 用户浏览过，则更新最后浏览时间，即updated_at。

			// 计算最后更新的时间是否超过viewAmountUpdateTime 6小时
			if diff := int(time.Now().Sub(targetView.ViewAmountUpdateTime).Hours()); diff > 6 { // 增加访问量，同时增加热度
				viewAdd = 1
				targetView.ViewAmountUpdateTime = time.Now()

				// 增加热度
				if ok := newsService.AddHotValueByNewsID(id, 1); !ok {
					response.ServerError(ctx, nil, "服务器错误")
				}

			} else {
				viewAdd = 0
			}
			// 无论是否增加，都要更新字段，可以更新UpdateAt字段
			mysqlDB.Model(&model.View{}).Updates(targetView)
		}

		update = bson.M{
			"$inc": bson.M{
				"view_amount": viewAdd,
			},
		}

		after := options.After
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}
		err = collection.FindOneAndUpdate(context.Background(), filter, update, &opt).Decode(&newsContent)
		if err != nil { // 此时找不到该新闻
			ctx.JSON(404, gin.H{
				"code": 404,
				"msg":  "Not found",
			})
		} else {
			newsContent.ConvertedTime = newsContent.OriginalTime.Format("2006年01月02日 15:04:05")
			b, _ := json.Marshal(&newsContent)
			var tmp gin.H
			err = json.Unmarshal(b, &tmp)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%T", tmp)
			ctx.JSON(200, gin.H{
				"code": 200,
				"data": tmp,
			})
		}
	} else { // 获取新闻列表
		var newsList []gin.H
		newsType := ctx.Query("news_type")
		fmt.Println(newsType)
		cur, err := collection.Find(context.Background(), bson.D{{"type", newsType}})

		if err != nil {
			log.Fatal(err)
		}
		for cur.Next(context.TODO()) {
			var news model.NormalNews
			var tmp gin.H
			err := cur.Decode(&news)
			if err != nil {
				log.Fatal(err)
			}
			b, _ := json.Marshal(&news)
			_ = json.Unmarshal(b, &tmp)
			newsList = append(newsList, tmp)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"news": newsList,
		})
	}
}

func GetHotNews(ctx *gin.Context) {
	db := common.InitRedis()
	defer db.Close()

	mongoDB, cancel := common.InitMongoDB()
	defer cancel()

	type dto struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	hotNewsList, _ := db.ZRevRange("hot_news", 0, 9).Result()

	var resultsDto []dto

	for _, v := range hotNewsList {
		id, _ := primitive.ObjectIDFromHex(v)

		var result model.NormalNews
		var resultDto dto
		filter := bson.M{"_id": id}
		_ = mongoDB.Collection("test").FindOne(context.Background(), filter).Decode(&result)

		resultDto.ID = v
		resultDto.Title = result.Title

		resultsDto = append(resultsDto, resultDto)
	}

	response.Success(ctx, gin.H{
		"hot_news": resultsDto,
	}, "")

	return
}
