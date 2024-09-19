package controller

import (
	"context"
	"log"
	"reflect"
	"strings"
	"time"
	"worldskills/models"
	"worldskills/resource"
	"worldskills/response"
	"worldskills/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var image_type_map = map[string]string{
	"title":       "string",
	"description": "string",
}

// 新增圖片
func (ctrl *Controller) PostImage(c *gin.Context) {
	var image models.Image

	access_token := c.GetHeader("X-Authorization")
	user := utils.Read(ctrl.Client, "worldskills", "users", bson.M{"access_token": access_token})
	if user == nil {
		response.Bad(c, "MSG_INVALID_ACCESS_TOKEN")
		return
	}

	// 檢查是否有缺少的欄位
	var check_keys = []string{"title", "description", "image"}
	if !utils.HasKey(c, check_keys) {
		response.Bad(c, "MSG_MISSING_FIELD")
		return
	}

	// 驗證資料型態
	for _, value := range check_keys {
		if value == "image" {
			continue
		}
		if reflect.TypeOf(c.PostForm(value)).String() != image_type_map[value] {
			response.Bad(c, "MSG_WRONG_DATA_TYPE")
			return
		}
	}

	// 儲存圖片
	image.UserID = user["_id"].(primitive.ObjectID)
	image.Title = c.PostForm("title")
	image.Description = c.PostForm("description")
	image.CreatedAt = utils.GetNowTime()
	image.ViewCount = 0

	// 圖片資訊
	file, _ := c.FormFile("image")

	// 檢查檔案格式
	mime, _ := utils.MimeType(file)
	if mime != "image/jpeg" && mime != "image/png" {
		response.Bad(c, "MSG_INVALID_FILE_FORMAT")
		return
	}

	image.MimeType = mime

	// 獲取圖片尺寸
	width, height, _ := utils.ImgSize(file)
	image.Width = width
	image.Height = height

	// 儲存檔案
	file_extension := strings.Split(file.Filename, ".")[1]
	file_name := time.Now().Format("20060102150405") + "." + file_extension
	path := "public/images/" + file_name
	c.SaveUploadedFile(file, path)

	// 儲存檔案路徑
	image.Url = "/" + path

	image_id := utils.Create(ctrl.Client, "worldskills", "images", image)
	data := utils.Read(ctrl.Client, "worldskills", "images", bson.M{"_id": image_id})
	resource.ImageDetailResource(c, &data, ctrl.Client)
	response.Ok(c, data)
}

// 軟刪除圖片
func (ctrl *Controller) DeleteImage(c *gin.Context) {
	access_token := c.GetHeader("X-Authorization")
	user := utils.Read(ctrl.Client, "worldskills", "users", bson.M{"access_token": access_token})

	// 檢查使用者是否存在
	if user == nil {
		response.Bad(c, "MSG_INVALID_ACCESS_TOKEN")
		return
	}

	// 檢查圖片是否存在
	image_id, _ := primitive.ObjectIDFromHex(c.Param("image_id"))
	image := utils.Read(ctrl.Client, "worldskills", "images", bson.M{"_id": image_id})
	if image == nil || image["deleted_at"] != "" {
		response.Bad(c, "MSG_IMAGE_NOT_EXISTS")
		return
	}

	// 檢查權限
	if image["user_id"] != user["_id"] && user["type"] != "ADMIN" {
		response.Bad(c, "MSG_PERMISSION_DENY")
		return
	}

	utils.Update(ctrl.Client, "worldskills", "images", bson.M{"_id": image_id}, bson.M{"$set": bson.M{"deleted_at": utils.GetNowTime()}})
	response.Ok(c, nil)
}

// 更新圖片
func (ctrl *Controller) PutImage(c *gin.Context) {
	access_token := c.GetHeader("X-Authorization")
	user := utils.Read(ctrl.Client, "worldskills", "users", bson.M{"access_token": access_token})

	request := utils.GetRequestData(c)
	// 檢查使用者是否存在
	if user == nil {
		response.Bad(c, "MSG_INVALID_ACCESS_TOKEN")
		return
	}
	log.Println(request)
	// 檢查圖片是否存在
	image_id, _ := primitive.ObjectIDFromHex(c.Param("image_id"))
	image := utils.Read(ctrl.Client, "worldskills", "images", bson.M{"_id": image_id})
	if image == nil || image["deleted_at"] != "" {
		response.Bad(c, "MSG_IMAGE_NOT_EXISTS")
		return
	}

	// 檢查權限
	if image["user_id"] != user["_id"] && user["type"] != "ADMIN" {
		response.Bad(c, "MSG_PERMISSION_DENY")
		return
	}

	// 驗證資料型態
	var pass_changes = []map[string]string{}
	for _, value := range []string{"title", "description"} {
		if request[value] == "" {
			continue
		}
		if reflect.TypeOf(request[value]).String() != image_type_map[value] {
			response.Bad(c, "MSG_WRONG_DATA_TYPE")
			return
		}
		pass_changes = append(pass_changes, map[string]string{value: request[value].(string)})
	}

	// 更新圖片資料
	change_data := bson.M{}
	for _, value := range pass_changes {
		for key, val := range value {
			change_data[key] = val
		}
	}
	change_data["updated_at"] = utils.GetNowTime()
	utils.Update(ctrl.Client, "worldskills", "images", bson.M{"_id": image_id}, bson.M{"$set": change_data})
	data := utils.Read(ctrl.Client, "worldskills", "images", bson.M{"_id": image_id})
	resource.ImageDetailResource(c, &data, ctrl.Client)
	response.Ok(c, data)
}

// 取得圖片
func (ctrl *Controller) GetImage(c *gin.Context) {
	image_id, _ := primitive.ObjectIDFromHex(c.Param("image_id"))
	data := utils.Read(ctrl.Client, "worldskills", "images", bson.M{"_id": image_id})

	// 檢查圖片是否存在
	if data == nil || data["deleted_at"] != "" {
		response.Bad(c, "MSG_IMAGE_NOT_EXISTS")
		return
	}

	// 更新圖片瀏覽次數
	utils.Update(ctrl.Client, "worldskills", "images", bson.M{"_id": image_id}, bson.M{"$inc": bson.M{"view_count": 1}})

	resource.ImageDetailResource(c, &data, ctrl.Client)
	response.Ok(c, data)
}

// 搜尋圖片
func (ctrl *Controller) Search(c *gin.Context) {
	// order_type: 1 -> asc, -1 -> desc
	order_type := map[string]int{
		"asc":  1,
		"desc": -1,
	}[c.Query("order_type")]
	if order_type == 0 {
		order_type = -1
	}

	// order_by: created_at, updated_at, width, height
	order_by_map := map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
		"width":      "width",
		"height":     "height",
	}

	var order_by string
	if order_by_map[c.Query("order_by")] == "" {
		order_by = "created_at"
	} else {
		order_by = order_by_map[c.Query("order_by")]
	}

	// page, page_size
	var page int
	var page_size int
	if page = utils.String2Int(c.Query("page")); page == 0 {
		page = 1
	}
	if page_size := utils.String2Int(c.Query("page_size")); page_size == 0 {
		page_size = 10
	}

	// 建立搜尋條件
	filter := bson.M{
		"deleted_at": "",
	}

	// 關鍵字搜尋(如果有)
	if keyword := c.Query("keyword"); keyword != "" {
		regex := bson.M{"$regex": ".*" + c.Query("keyword") + ".*"}
		filter["$or"] = []bson.M{
			{"title": regex},
			{"description": regex},
		}
	}

	// 搜尋圖片
	collection := ctrl.Client.Database("worldskills").Collection("images")
	collection.Find(context.TODO(), filter)

	// 查詢的Options
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: order_by, Value: order_type}})
	findOptions.SetSkip(int64((page - 1) * page_size))
	limit := int64(page_size)
	limit64 := int64(limit)
	findOptions.Limit = &limit64
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 解碼結果
	var results []map[string]interface{}
	for cursor.Next(context.Background()) {
		var result map[string]interface{}
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		resource.ImageResource(c, &result)
		results = append(results, result)
	}

	response.Ok(c, results)
}

// 取得受歡迎的圖片
func (ctrl *Controller) GetPopularImages(c *gin.Context) {
	var limit = utils.String2Int(c.Query("limit"))
	if limit < 0 || limit > 100 {
		response.Bad(c, "MSG_WRONG_DATA_TYPE")
		return
	}

	if limit == 0 {
		limit = 10
	}

	// 建立搜尋條件
	filter := bson.M{
		"deleted_at": "",
	}

	// 限制搜尋筆數
	findOptions := options.Find()
	limit64 := int64(limit)
	findOptions.Limit = &limit64

	// 搜尋圖片
	results := utils.ReadAll(ctrl.Client, "worldskills", "images", filter, findOptions)

	response.Ok(c, results)
}
