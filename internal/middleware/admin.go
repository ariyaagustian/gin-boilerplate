package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	return func(c *gin.Context) {
		if adminEmail == "" {
			c.Next() // skip check kalau tidak diset
			return
		}
		email, _ := c.Get("user_email")

		if email != adminEmail {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
