package commentService

import (
	"context"
	"gin_vue_project/common"
	"gin_vue_project/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func GetCommentAmount(newsID string) int {
	db, cancel := common.InitMongoDB()
	defer cancel()

	newsObjectID, _ := primitive.ObjectIDFromHex(newsID)
	filter := bson.M{
		"_id": newsObjectID,
	}

	opt := options.FindOneOptions{
		Projection: bson.D{
			{"_id", 0},
			{"comment_amount", 1},
		},
	}
	//var amount int
	var normalNews model.NormalNews
	err := db.Collection("test").FindOne(context.TODO(), filter, &opt).Decode(&normalNews)

	if err != nil {
		log.Println(err.Error())
		return -1
	}
	return normalNews.CommentAmount
}
