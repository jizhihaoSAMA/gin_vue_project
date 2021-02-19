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
	ID            string    `json:"id" bson:"_id"`
	URL           string    `json:"url" bson:"url"`
	Title         string    `json:"title" bson:"title"`
	Content       bson.A    `json:"content"`
	From          string    `json:"from"`
	Images        bson.A    `json:"images"`
	Type          string    `json:"type" bson:"type"`
	CommentAmount int       `json:"comment_amount" bson:"comment_amount"`
	ViewAmount    int       `json:"view_amount" bson:"view_amount"`
	Layout        string    `json:"layout" bson:"layout"`
	OriginalTime  time.Time `bson:"time" json:"-"`
	ConvertedTime string    `json:"time"`
}
