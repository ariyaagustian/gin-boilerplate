package http

import (
	"github.com/gin-gonic/gin"

	"github.com/ariyaagustian/gin-crud-boilerplate/internal/middleware"
	"github.com/ariyaagustian/gin-crud-boilerplate/internal/transport/http/handler"
)

// internal/transport/http/router.go
func NewRouter(userH *handler.UserHandler, authH *handler.AuthHandler, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), middleware.CORS())

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)

	api := r.Group("/api/v1", middleware.AuthBearer(jwtSecret))
	{
		// admin only
		admin := api.Group("/admin", middleware.AdminOnly())
		{
			admin.POST("/users/set-password", authH.AdminSetPassword)
		}

		u := api.Group("/users")
		{
			u.POST("", userH.Create)
			u.GET("", userH.List)
			u.GET("/:id", userH.Get)
			u.PUT("/:id", userH.Update)
			u.DELETE("/:id", userH.Delete)
			u.GET("/me", userH.Me)
		}
	}
	return r
}
