package authinfrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appauth "pethelp-backend/internal/app/auth"
)

type RefreshDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshHandler returns a Gin handler that processes refresh requests.
// It reads the old refresh token, calls AuthService.RefreshToken, and returns new tokens.
func RefreshHandler(
	svc *appauth.AuthService,
	logger *zap.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto RefreshDTO
		// Bind and validate JSON payload
		if err := c.ShouldBindJSON(&dto); err != nil {
			logger.Warn("invalid refresh payload", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Perform token refresh
		accessToken, refreshToken, err := svc.RefreshToken(c.Request.Context(), dto.RefreshToken)
		if err != nil {
			logger.Warn("refresh failed", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
			return
		}

		// Return new tokens
		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}
