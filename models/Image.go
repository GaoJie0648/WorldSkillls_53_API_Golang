package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Image struct {
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Url         string             `json:"url" bson:"url"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Width       int                `json:"width" bson:"width"`
	Height      int                `json:"height" bson:"height"`
	MimeType    string             `json:"mimetype" bson:"mimetype"`
	ViewCount   int                `json:"view_count" bson:"view_count"`
	CreatedAt   string             `json:"created_at" bson:"created_at"`
	UpdatedAt   string             `json:"updated_at" bson:"updated_at"`
	DeletedAt   string             `json:"deleted_at" bson:"deleted_at"`
}
