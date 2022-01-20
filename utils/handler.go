package utils

import (
	"github.com/gin-gonic/gin"
)

// Response will return json responce of http
// This function hsould handle both error and sucess
func Response(c *gin.Context, statusCode int, payload interface{}) {
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(statusCode, payload)
}

func SuccessPayload(data interface{}, message string, args ...int) map[string]interface{} {
	result := make(map[string]interface{})
	result["data"] = data
	result["message"] = message
	return result
}
