package response

import (
	"net/http"

	"github.com/ariyaagustian/gin-boilerplate/pkg/apperr"
	"github.com/gin-gonic/gin"
)

// WriteError mengirimkan error JSON konsisten
func WriteError(c *gin.Context, err error) {
	if ae, ok := err.(*apperr.AppError); ok {
		// kalau 5xx, log otomatis
		if ae.HTTPStatus >= 500 {
			c.Error(ae)
		}
		c.JSON(ae.HTTPStatus, gin.H{
			"error": gin.H{
				"code":    ae.Code,
				"message": ae.Message,
			},
		})
		return
	}

	// fallback: error unknown
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":    "internal",
			"message": "terjadi kesalahan pada server",
		},
	})
}
