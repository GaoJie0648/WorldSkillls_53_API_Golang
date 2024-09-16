package resource

import (
	"worldskills/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ImageDetailResource(c *gin.Context, data *map[string]interface{}, client *mongo.Client) {
	var tmp = *data
	// 取得使用者資料
	user := utils.Read(client, "worldskills", "users", bson.M{"_id": tmp["user_id"]})
	UserResource(c, &user)

	// 取得圖片評論數
	comment_count := utils.ReadAll(client, "worldskills", "comments", bson.M{"image_id": tmp["_id"]}, nil)

	response := map[string]interface{}{
		"id":            tmp["_id"],
		"url":           tmp["url"],
		"author":        user,
		"title":         tmp["title"],
		"description":   tmp["description"],
		"width":         tmp["width"],
		"height":        tmp["height"],
		"mimetype":      tmp["mimetype"],
		"view_count":    tmp["view_count"],
		"comment_count": len(comment_count),
		"created_at":    tmp["created_at"],
		"updated_at":    tmp["updated_at"],
	}

	*data = response
}
