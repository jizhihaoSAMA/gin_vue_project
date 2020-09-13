package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type News struct {
	Title string    `json:"title" bson:"title"`
	Time  time.Time `json:"time" bson:"time"`
	ID    string    `json:"id" bson:"_id"`
}

type NormalNews struct {
	Title   string `json:"title" bson:"title"`
	Content bson.A `json:"content"`
	From    string `json:"from"`
	Images  bson.A `json:"images"`
	Type    string `json:"type" bson:"type"`
}
