package controller

import (
	"encoding/json"
	"io"
	"reflect"
	"worldskills/models"
	"worldskills/resource"
	"worldskills/response"
	"worldskills/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var comment_type_map = map[string]string{
	"content": "string",
}

// 新增評論
func (ctrl *Controller) PostComment(c *gin.Context) {
	var comment models.Comment
	image_id, _ := primitive.ObjectIDFromHex(c.Param("image_id"))
	// 讀取 Request Body 內容
	body, _ := io.ReadAll(c.Request.Body)
	var bodyMap map[string]interface{}
	json.Unmarshal(body, &bodyMap)

	// 檢查使用者是否存在
	access_token := c.GetHeader("X-Authorization")
	user := utils.Read(ctrl.Client, "worldskills", "users", bson.M{"access_token": access_token})
	if user == nil {
		response.Bad(c, "MSG_INVALID_ACCESS_TOKEN")
		return
	}

	// 檢查是否有缺少的欄位
	var check_keys = []string{"content"}
	if bodyMap["content"] == nil {
		response.Bad(c, "MSG_MISSING_FIELD")
		return
	}

	// 驗證資料型態
	for _, value := range check_keys {
		if reflect.TypeOf(bodyMap[value]).String() != comment_type_map[value] {
			response.Bad(c, "MSG_WRONG_DATA_TYPE")
			return
		}
	}

	// 檢查圖片是否存在
	image := utils.Read(ctrl.Client, "worldskills", "images", bson.M{"_id": image_id})
	if image == nil || image["deleted_at"] != "" {
		response.Bad(c, "MSG_IMAGE_NOT_EXISTS")
		return
	}

	// 新增評論
	comment.UserID = user["_id"].(primitive.ObjectID)
	comment.ImageID = image["_id"].(primitive.ObjectID)
	comment.Content = bodyMap["content"].(string)
	comment.CreatedAt = utils.GetNowTime()

	comment_id := utils.Create(ctrl.Client, "worldskills", "comments", comment)
	data := utils.Read(ctrl.Client, "worldskills", "comments", bson.M{"_id": comment_id})
	resource.CommentResource(c, &data, ctrl.Client)
	response.Ok(c, data)
}

// 刪除評論
func (ctrl *Controller) DeleteComment(c *gin.Context) {
	comment_id, _ := primitive.ObjectIDFromHex(c.Param("comment_id"))
	image_id, _ := primitive.ObjectIDFromHex(c.Param("image_id"))

	// 檢查使用者是否存在
	access_token := c.GetHeader("X-Authorization")
	user := utils.Read(ctrl.Client, "worldskills", "users", bson.M{"access_token": access_token})
	if user == nil {
		response.Bad(c, "MSG_INVALID_ACCESS_TOKEN")
		return
	}

	// 檢查評論是否存在
	comment := utils.Read(ctrl.Client, "worldskills", "comments", bson.M{"_id": comment_id})
	if comment == nil {
		response.Bad(c, "MSG_COMMENT_NOT_EXISTS")
		return
	}

	// 檢查圖片是否存在
	image := utils.Read(ctrl.Client, "worldskills", "images", bson.M{"_id": image_id})
	if image == nil || image["deleted_at"] != "" {
		response.Bad(c, "MSG_IMAGE_NOT_EXISTS")
		return
	}

	// 檢查使用者是否有權限刪除評論
	if comment["user_id"] != user["_id"] && user["type"] != "ADMIN" {
		response.Bad(c, "MSG_PERMISSION_DENY")
		return
	}

	// 刪除評論
	utils.Delete(ctrl.Client, "worldskills", "comments", bson.M{"_id": comment_id})
	response.Ok(c, nil)
}

// 取得圖片評論
func (ctrl *Controller) GetImageComments(c *gin.Context) {
	image_id, _ := primitive.ObjectIDFromHex(c.Param("image_id"))

	// 檢查圖片是否存在
	image := utils.Read(ctrl.Client, "worldskills", "images", bson.M{"_id": image_id})
	if image == nil || image["deleted_at"] != "" {
		response.Bad(c, "MSG_IMAGE_NOT_EXISTS")
		return
	}

	// 取得評論
	comments := utils.ReadAll(ctrl.Client, "worldskills", "comments", bson.M{"image_id": image_id})
	comments_data := []map[string]interface{}{}
	for _, comment := range comments {
		resource.CommentResource(c, &comment, ctrl.Client)
		comments_data = append(comments_data, comment)
	}
	response.Ok(c, comments_data)
}

// 回覆評論
func (ctrl *Controller) ReplyComment(c *gin.Context) {
	var comment models.Comment
	comment_id, _ := primitive.ObjectIDFromHex(c.Param("comment_id"))
	image_id, _ := primitive.ObjectIDFromHex(c.Param("image_id"))

	// 讀取 Request Body 內容
	body, _ := io.ReadAll(c.Request.Body)
	var bodyMap map[string]interface{}
	json.Unmarshal(body, &bodyMap)

	// 檢查使用者是否存在
	access_token := c.GetHeader("X-Authorization")
	user := utils.Read(ctrl.Client, "worldskills", "users", bson.M{"access_token": access_token})
	if user == nil {
		response.Bad(c, "MSG_INVALID_ACCESS_TOKEN")
		return
	}

	// 檢查圖片是否存在
	image := utils.Read(ctrl.Client, "worldskills", "images", bson.M{"_id": image_id})
	if image == nil || image["deleted_at"] != "" {
		response.Bad(c, "MSG_IMAGE_NOT_EXISTS")
		return
	}

	// 檢查評論是否存在
	comments_data := utils.Read(ctrl.Client, "worldskills", "comments", bson.M{"_id": comment_id})
	if comments_data == nil {
		response.Bad(c, "MSG_COMMENT_NOT_EXISTS")
		return
	}

	// 檢查是否有缺少的欄位
	var check_keys = []string{"content"}
	if bodyMap["content"] == nil {
		response.Bad(c, "MSG_MISSING_FIELD")
		return
	}

	// 驗證資料型態
	for _, value := range check_keys {
		if reflect.TypeOf(bodyMap[value]).String() != comment_type_map[value] {
			response.Bad(c, "MSG_WRONG_DATA_TYPE")
			return
		}
	}

	// 新增回覆
	comment.UserID = user["_id"].(primitive.ObjectID)
	comment.ImageID = image["_id"].(primitive.ObjectID)
	comment.Content = bodyMap["content"].(string)
	comment.CreatedAt = utils.GetNowTime()
	comment.ReplyID = comment_id

	created_id := utils.Create(ctrl.Client, "worldskills", "comments", comment)
	data := utils.Read(ctrl.Client, "worldskills", "comments", bson.M{"_id": created_id})
	resource.CommentResource(c, &data, ctrl.Client)
	response.Ok(c, data)
}
