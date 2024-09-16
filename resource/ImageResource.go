package resource

import "github.com/gin-gonic/gin"

func ImageResource(c *gin.Context, data *map[string]interface{}) {
	var tmp = *data
	response := map[string]interface{}{
		"id":          tmp["_id"],
		"url":         tmp["url"],
		"title":       tmp["title"],
		"description": tmp["description"],
		"created_at":  tmp["created_at"],
		"updated_at":  tmp["updated_at"],
	}

	*data = response
}
