package response

import "github.com/gin-gonic/gin"

func ResponseOK(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func ResponseNOK(c *gin.Context, code int, message string, errors interface{}) {
	c.JSON(code, gin.H{
		"success": false,
		"message": message,
		"errors":  errors,
	})
}
