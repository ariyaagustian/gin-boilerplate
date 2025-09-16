package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ariyaagustian/gin-boilerplate/pkg/auth"
)

func AuthBearer(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		hdr := c.GetHeader("Authorization")
		if !strings.HasPrefix(hdr, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimPrefix(hdr, "Bearer ")
		claims, err := auth.Parse(token, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		uid, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid uid"})
			return
		}
		// inject ke context
		c.Set("user_id", uid)
		if claims.Email != "" {
			c.Set("user_email", claims.Email)
		}
		c.Next()
	}
}
