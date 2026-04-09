package utils

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code":    200,
		"success": true,
		"data":    data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"success": false,
		"message": message,
	})
}
