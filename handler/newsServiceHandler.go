package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_vue_project/common"
	"gin_vue_project/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func GetNews(ctx *gin.Context) {
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
		if true {
			update = bson.M{
				"$inc": bson.M{
					"view_amount": 1,
				},
			}
		} else {
			update = bson.M{
				"$inc": bson.M{
					"view_amount": 0,
				},
			}
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
		return
	}
	// 获取列表
	var newsList []gin.H
	newsType := ctx.Query("news_type")
	fmt.Println(newsType)
	cur, err := collection.Find(context.Background(), bson.D{{"type", newsType}})

	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var news model.News
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
