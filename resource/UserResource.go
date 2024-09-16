package resource

import (
	"github.com/gin-gonic/gin"
)

func UserResource(c *gin.Context, data *map[string]interface{}) {
	var tmp = *data
	response := map[string]interface{}{
		"id":            tmp["_id"],
		"email":         tmp["email"],
		"nickname":      tmp["nickname"],
		"profile_image": tmp["profile_image"],
		"type":          tmp["type"],
		"created_at":    tmp["created_at"],
	}
	if c.Request.URL.Path == "/api/auth/login" {
		response["access_token"] = tmp["access_token"]
	}

	*data = response
}
