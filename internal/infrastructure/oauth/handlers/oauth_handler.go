package handlers

import (
	"fmt"
	"net/http"

	oauthEnt "pethelp-backend/internal/domain/oauth"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

// OAuthHandlers OAuth2.0 services contains
type OAuthHandlers struct {
	OAuthService oauthEnt.OAuthTokenService
	Logger       *zap.Logger
}

const operationName = "oauth_token_service:"

// NewOAuthHandlers create new OAuthHandlers
func NewOAuthHandlers(service oauthEnt.OAuthTokenService, logger *zap.Logger) *OAuthHandlers {
	return &OAuthHandlers{OAuthService: service, Logger: logger}
}

// SignInWithProvider redirect to provider login page
func (h *OAuthHandlers) SignInWithProvider(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// CallbackHandler process provider callback and save user data
func (h *OAuthHandlers) ProviderCallback(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		messErr := fmt.Errorf("%s failed to complete OAuth2.0 authentication: %w", operationName, err)
		h.Logger.Error("", zap.Error(messErr))

		errMessage := oauthEnt.OAuthErrResponse{
			Code:    http.StatusInternalServerError,
			Type:    "OAuth error",
			Message: messErr.Error(),
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errMessage)
		return
	}

	err = h.OAuthService.InitAuth(c.Request.Context(), &user)
	if err != nil {
		messErr := fmt.Errorf("%s failed to save auth data: %w", operationName, err)
		h.Logger.Error("", zap.Error(messErr))

		errMessage := oauthEnt.OAuthErrResponse{
			Code:    http.StatusInternalServerError,
			Type:    "DB error",
			Message: messErr.Error(),
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, errMessage)
		return
	}

	tokenData := oauthEnt.TokensData{
		Access:  user.AccessToken,
		Refresh: user.RefreshToken,
		ID:      user.IDToken,
	}

	// — Success response
	c.JSON(http.StatusCreated, tokenData)

}
