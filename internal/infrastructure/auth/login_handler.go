package authinfrastructure

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appauth "pethelp-backend/internal/app/auth"
	dom "pethelp-backend/internal/domain/auth"
)

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(svc *appauth.AuthService, logger *zap.Logger)  gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto LoginDTO

		//bind and validate JSON payload
		if err := c.ShouldBindJSON(&dto); err != nil {
			logger.Warn("invalid login payload", zap.String("error", "validation failed"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//perform authentication and token issuance
		accessToken, refreshToken, err := svc.Login(c.Request.Context(), dto.Email, dto.Password)
		if err != nil {
			if errors.Is(err, dom.ErrInvalidCredentials) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			} else {
				logger.Error("login failed", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
			return
		}

		//return both tokens in the response body
		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}

}