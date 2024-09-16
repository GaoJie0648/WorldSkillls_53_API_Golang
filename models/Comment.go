package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	ImageID   primitive.ObjectID `json:"image_id" bson:"image_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Content   string             `json:"content" bson:"content"`
	ReplyID   primitive.ObjectID `json:"reply_id" bson:"reply_id"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
}
