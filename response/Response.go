package response

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Success bool

var codes = map[string]int{
	"MSG_INVALID_LOGIN":        403,
	"MSG_USER_EXISTS":          409,
	"MSG_PASSWORD_NOT_SECURE":  409,
	"MSG_INVALID_ACCESS_TOKEN": 401,
	"MSG_PERMISSION_DENY":      403,
	"MSG_MISSING_FIELD":        400,
	"MSG_WRONG_DATA_TYPE":      400,
	"MSG_IMAGE_NOT_EXISTS":     404,
	"MSG_COMMENT_NOT_EXISTS":   404,
	"MSG_USER_NOT_EXISTS":      404,
	"MSG_INVALID_FILE_FORMAT":  400,
}

type Accept struct {
	Success `json:"success"`
	Data    interface{} `json:"data"`
}

type Reject struct {
	Success `json:"success"`
	Message string `json:"message"`
}

func Ok(c *gin.Context, data interface{}) {
	if data == nil {
		c.JSON(http.StatusOK, struct {
			Success bool `json:"success"`
		}{
			Success: true,
		})
	} else {
		c.JSON(http.StatusOK, Accept{
			Success: true,
			Data:    data,
		})
	}
}

func Bad(c *gin.Context, message string) {
	if status := codes[message]; status != 0 {
		c.JSON(status, Reject{
			Success: false,
			Message: message,
		})
	} else {
		log.Fatal("Invalid message code")
	}
}
