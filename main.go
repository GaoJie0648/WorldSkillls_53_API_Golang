package main

import (
	"context"
	"net/http"
	"time"
	"worldskills/controller"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

// 初始化資料庫
func initDB() {
	var err error
	var uri = "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
}

func main() {

	initDB()

	router := gin.Default()

	router.StaticFS("/public", http.Dir("public"))

	ctrl := controller.GetClient(client)

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", ctrl.Register)
			auth.POST("/login", ctrl.Login)
			auth.POST("/logout", ctrl.Logout)
		}

		image := api.Group("/image")
		{
			image.POST("/upload", ctrl.PostImage)
			image.DELETE("/:image_id", ctrl.DeleteImage)
			image.PUT("/:image_id", ctrl.PutImage)
			image.GET("/:image_id", ctrl.GetImage)

			image.POST("/:image_id/comment", ctrl.PostComment)
			image.DELETE("/:image_id/comment/:comment_id", ctrl.DeleteComment)
			image.GET("/:image_id/comment", ctrl.GetImageComments)

			image.POST("/:image_id/comment/:comment_id/reply", ctrl.ReplyComment)

			image.GET("/search", ctrl.Search)
		}

		user := api.Group("/user")
		{
			user.GET("/:user_id/image", ctrl.GetUserImages)
		}
	}

	router.Run(":8080")
}
