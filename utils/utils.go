package utils

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
)

func ImgSize(file *multipart.FileHeader) (width int, height int, err error) {
	src, err := file.Open()
	if err != nil {
		return 0, 0, err
	}
	defer src.Close()

	// 解碼圖像
	img, _, err := image.Decode(src)
	if err != nil {
		return 0, 0, err
	}

	// 獲取圖像的寬度和高度
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	return w, h, nil
}

func MimeType(file *multipart.FileHeader) (string, error) {

	src, _ := file.Open()
	defer src.Close()

	buffer := make([]byte, 512)
	_, err := src.Read(buffer)
	if err != nil {
		return "", err
	}

	mime := mimetype.Detect(buffer)
	return mime.String(), nil
}

func HasKey(c *gin.Context, key []string) bool {
	for _, k := range key {
		file, _ := c.FormFile(k)
		if c.PostForm(k) == "" && file == nil {
			return false
		}
	}
	return true
}

func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func String2Int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil || i <= 0 {
		i = 10
	}
	return int(i)
}
