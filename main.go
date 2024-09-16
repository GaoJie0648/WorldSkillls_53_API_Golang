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
			auth.POST("/login", ctrl.Login)       // API 1
			auth.POST("/register", ctrl.Register) // API 2
			auth.POST("/logout", ctrl.Logout)     // API 3
		}

		image := api.Group("/image")
		{
			image.GET("/search", ctrl.Search)            // API 4
			image.POST("/upload", ctrl.PostImage)        // API 7
			image.PUT("/:image_id", ctrl.PutImage)       // API 8
			image.GET("/:image_id", ctrl.GetImage)       // API 9
			image.DELETE("/:image_id", ctrl.DeleteImage) // API 10

			image.GET("/:image_id/comment", ctrl.GetImageComments)                // API 11
			image.POST("/:image_id/comment", ctrl.PostComment)                    // API 12
			image.POST("/:image_id/comment/:comment_id/reply", ctrl.ReplyComment) // API 13
			image.DELETE("/:image_id/comment/:comment_id", ctrl.DeleteComment)    // API 14

		}

		user := api.Group("/user")
		{
			user.GET("/:user_id/image", ctrl.GetUserImages) // API 6
		}
	}

	router.Run(":8080")
}
