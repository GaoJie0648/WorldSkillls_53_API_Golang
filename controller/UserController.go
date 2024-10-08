package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"worldskills/models"
	"worldskills/resource"
	"worldskills/response"
	"worldskills/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var user_type_map = map[string]string{
	"email":    "string",
	"nickname": "string",
	"password": "string",
}

// 登入
func (ctrl *Controller) Login(c *gin.Context) {
	// 檢查是否有缺少的欄位
	var check_keys = []string{"email", "password"}
	request := utils.GetRequestData(c)
	for _, key := range check_keys {
		if request[key] == nil {
			response.Bad(c, "MSG_MISSING_FIELD")
			return
		}
	}

	// 檢查使用者是否存在
	var user_data = utils.Read(ctrl.Client, "worldskills", "users", bson.M{"email": request["email"]})
	if user_data == nil {
		response.Bad(c, "MSG_USER_NOT_EXISTS")
		return
	}

	// 檢查密碼是否正確
	if err := bcrypt.CompareHashAndPassword([]byte(user_data["password"].(string)), []byte(request["password"].(string))); err != nil {
		response.Bad(c, "MSG_INVALID_LOGIN")
		return
	}

	// 產生 Access Token
	hasher := sha256.New()
	hasher.Write([]byte(user_data["email"].(string)))
	user_data["access_token"] = hex.EncodeToString(hasher.Sum(nil))
	// 更新使用者資料
	utils.Update(ctrl.Client, "worldskills", "users", bson.M{"email": request["email"]}, bson.M{"$set": bson.M{"access_token": user_data["access_token"]}})

	resource.UserResource(c, &user_data)
	response.Ok(c, user_data)
}

// 登出
func (ctrl *Controller) Logout(c *gin.Context) {
	// 檢查是否有缺少的欄位
	var check_keys = []string{"email", "password"}
	if !utils.HasKey(c, check_keys) {
		response.Bad(c, "MSG_MISSING_FIELD")
		return
	}

	// 檢查使用者是否存在
	var user_data = utils.Read(ctrl.Client, "worldskills", "users", bson.M{"email": c.PostForm("email")})
	if user_data == nil {
		response.Bad(c, "MSG_USER_NOT_EXISTS")
		return
	}

	// 檢查密碼是否正確
	if err := bcrypt.CompareHashAndPassword([]byte(user_data["password"].(string)), []byte(c.PostForm("password"))); err != nil {
		response.Bad(c, "MSG_INVALID_LOGIN")
		return
	}

	// 刪除 Access Token
	utils.Update(ctrl.Client, "worldskills", "users", bson.M{"email": c.PostForm("email")}, bson.M{"$unset": bson.M{"access_token": nil}})
	response.Ok(c, nil)
}

// 新增使用者
func (ctrl *Controller) Register(c *gin.Context) {
	var user models.User

	// 檢查是否有缺少的欄位
	var check_keys = []string{"email", "nickname", "password", "profile_image"}
	if !utils.HasKey(c, check_keys) {
		response.Bad(c, "MSG_MISSING_FIELD")
		return
	}

	// 驗證電子郵件格式
	if !strings.Contains(c.PostForm("email"), "@") {
		response.Bad(c, "MSG_WRONG_DATA_TYPE")
		return
	}

	// 驗證資料型態
	for _, value := range check_keys {
		if value == "profile_image" {
			continue
		}
		if reflect.TypeOf(c.PostForm(value)).String() != user_type_map[value] {
			response.Bad(c, "MSG_WRONG_DATA_TYPE")
			return
		}
	}

	// 檢查使用者是否已經存在
	var user_check = utils.Read(ctrl.Client, "worldskills", "users", bson.M{"email": c.PostForm("email")})
	if user_check != nil {
		response.Bad(c, "MSG_USER_EXISTS")
		return
	}

	// 檢查密碼是否安全
	if len(c.PostForm("password")) <= 4 {
		response.Bad(c, "MSG_PASSWORD_NOT_SECURE")
		return
	}

	// 儲存使用者資料

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), bcrypt.DefaultCost)
	user.Email = c.PostForm("email")
	user.Nickname = c.PostForm("nickname")
	user.Password = string(hashedPassword)
	user.Type = "USER"
	user.CreatedAt = utils.GetNowTime()

	// 使用者頭像資料
	file, _ := c.FormFile("profile_image")

	// 檢查檔案格式
	mime, _ := utils.MimeType(file)
	if mime != "image/jpeg" && mime != "image/png" {
		response.Bad(c, "MSG_INVALID_FILE_FORMAT")
		return
	}

	// 儲存檔案
	file_extension := strings.Split(file.Filename, ".")[1]
	file_name := time.Now().Format("20060102150405") + "." + file_extension
	path := "public/profile_images/" + file_name
	c.SaveUploadedFile(file, path)

	// 儲存檔案路徑
	user.ProfileImage = "/" + path

	user_id := utils.Create(ctrl.Client, "worldskills", "users", user)
	data := utils.Read(ctrl.Client, "worldskills", "users", bson.M{"_id": user_id})
	resource.UserResource(c, &data)
	response.Ok(c, data)
}

// 取得使用者圖片
func (ctrl *Controller) GetUserImages(c *gin.Context) {
	user_id, _ := primitive.ObjectIDFromHex(c.Param("user_id"))
	user := utils.Read(ctrl.Client, "worldskills", "users", bson.M{"_id": user_id})

	// 檢查使用者是否存在
	if user == nil {
		response.Bad(c, "MSG_USER_NOT_EXISTS")
		return
	}

	images := utils.ReadAll(ctrl.Client, "worldskills", "images", bson.M{"user_id": user_id, "deleted_at": ""}, nil)
	images_data := []map[string]interface{}{}

	for _, image := range images {
		resource.ImageResource(c, &image)
		images_data = append(images_data, image)
	}

	response.Ok(c, images_data)
}

// 取得受歡迎的使用者
func (ctrl *Controller) GetPopularUsers(c *gin.Context) {
	users := utils.ReadAll(ctrl.Client, "worldskills", "users", bson.M{}, nil)
	for _, user := range users {
		images := utils.ReadAll(ctrl.Client, "worldskills", "images", bson.M{"user_id": user["_id"], "deleted_at": ""}, nil)
		user["image_count"] = len(images)
		user["upload_count"] = len(utils.ReadAll(ctrl.Client, "worldskills", "images", bson.M{"user_id": user["_id"], "deleted_at": "", "updated_at": bson.M{"$ne": ""}}, nil))
		comments := 0
		for _, image := range images {
			image_id := image["_id"].(primitive.ObjectID)
			comments += len(utils.ReadAll(ctrl.Client, "worldskills", "comments", bson.M{"image_id": image_id}, nil))
		}
		user["total_comment_count"] = comments
	}

	order_by := c.Query("order_by")
	if order_by == "" {
		order_by = "upload_count"
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i][order_by].(int) > users[j][order_by].(int)
	})

	limit_str := c.Query("limit")
	limit, err := strconv.Atoi(limit_str)
	if err != nil || limit <= 0 {
		limit = 10
	}

	if len(users) > limit {
		users = users[:limit]
	}

	data_map := []map[string]interface{}{}
	for _, user := range users {
		tmp := user
		resource.UserResource(c, &user)
		data := map[string]interface{}{
			"user":                user,
			"image_count":         tmp["image_count"],
			"upload_count":        tmp["upload_count"],
			"total_comment_count": tmp["total_comment_count"],
		}
		data_map = append(data_map, data)
	}
	response.Ok(c, data_map)
}
