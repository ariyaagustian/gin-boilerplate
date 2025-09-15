package response

import "github.com/gin-gonic/gin"

func JSON(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{"data": data})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}
