package utils

import (
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"status":  code,
		"message": message,
		"data":    data,
	})
}
