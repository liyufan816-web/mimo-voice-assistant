package utils

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"success": true,
		"data":    data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"message": message,
	})
}
