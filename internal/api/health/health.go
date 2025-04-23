package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Check returns http handler for server health check
func Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	}
}