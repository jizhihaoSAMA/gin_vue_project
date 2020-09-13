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
	"log"
	"net/http"
)

func GetNews(ctx *gin.Context) {
	db, cancel := common.InitMongoDB()
	defer cancel()
	collection := db.Collection("test")

	id := ctx.Query("id")
	if id != "" {
		var newsContent model.NormalNews
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Fatal(err)
		}
		filter := bson.D{{"_id", objectID}}
		err = collection.FindOne(context.Background(), filter).Decode(&newsContent)
		if err != nil {
			fmt.Println(err.Error())
			ctx.JSON(404, gin.H{
				"code": 404,
				"msg":  "Not found",
			})
		} else {
			b, _ := json.Marshal(&newsContent)
			var tmp gin.H
			err = json.Unmarshal(b, &tmp)
			if err != nil {
				log.Fatal(err)
			}
			ctx.JSON(200, gin.H{
				"code": 200,
				"data": tmp,
			})
		}
		return
	}

	var newsList []gin.H

	cur, err := collection.Find(context.Background(), bson.D{{}})

	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var news model.News
		var tmp gin.H
		err := cur.Decode(&news)
		fmt.Println(news)
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
