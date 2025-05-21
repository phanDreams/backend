package oauth

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
)

type TokensData struct {
	Access  string `json:"access_token,omitempty"`
	Refresh string `json:"refresh_token,omitempty"`
	ID      string `json:"id_token,omitempty"`
}

type OAuthErrResponse struct {
	Code    int    `json:"code,omitempty"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}

type OAuthTokenService interface {
	InitAuth(ctx context.Context, user *goth.User) error
	VerifyAuth(ctx context.Context, identifier string) (*goth.User, error)
	RefreshAuth(ctx context.Context, identifier string) (string, error)
	RevokeAuth(ctx context.Context, c *gin.Context, identifier string) error
}

type OAuthTokenRepository interface {
	SetToken(ctx context.Context, user *goth.User) error
	GetToken(ctx context.Context, identifier string) (*goth.User, error)
	DelToken(ctx context.Context, identifier string) error
}
