package response

import "github.com/gin-gonic/gin"

func JSON(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"success": true, "data": data})
}

func Error(c *gin.Context, err error, code int) {
	c.JSON(code, gin.H{"success": false, "error": gin.H{"code": code, "message": err.Error()}})
}
