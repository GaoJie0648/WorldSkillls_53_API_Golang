package controller

import "go.mongodb.org/mongo-driver/mongo"

type Controller struct {
	Client *mongo.Client
}

func GetClient(client *mongo.Client) *Controller {
	return &Controller{Client: client}
}
