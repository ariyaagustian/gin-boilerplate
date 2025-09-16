package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ariyaagustian/gin-boilerplate/internal/middleware"
	"github.com/ariyaagustian/gin-boilerplate/internal/transport/http/handler"
	"gorm.io/gorm"
)

// internal/transport/http/router.go
func NewRouter(userH *handler.UserHandler, authH *handler.AuthHandler, jwtSecret string, db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), middleware.CORS())

	// Health check endpoints
	healthH := handler.NewHealthHandler(db)
	r.GET("/healthz", healthH.HealthCheck)
	r.GET("/health/liveness", healthH.Liveness)
	r.GET("/health/readiness", healthH.Readiness)

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

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
