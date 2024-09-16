package resource

import (
	"worldskills/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CommentResource(c *gin.Context, data *map[string]interface{}, client *mongo.Client) {
	var tmp = *data
	// 取得使用者資料
	user := utils.Read(client, "worldskills", "users", bson.M{"_id": tmp["user_id"]})
	UserResource(c, &user)

	// 取得評論回覆數
	replys := utils.ReadAll(client, "worldskills", "comments", bson.M{"reply_id": tmp["_id"]}, nil)
	comments := []map[string]interface{}{}
	for _, reply := range replys {
		CommentResource(c, &reply, client)
		comments = append(comments, reply)
	}

	response := map[string]interface{}{
		"id":         tmp["_id"],
		"user":       user,
		"content":    tmp["content"],
		"created_at": tmp["created_at"],
		"comments":   comments,
	}

	*data = response
}
