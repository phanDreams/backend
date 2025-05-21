package service

import (
	"context"
	"fmt"

	"pethelp-backend/internal/config"
	oauthEnt "pethelp-backend/internal/domain/oauth"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const operationName = "oauth_token_service:"

type OAuthServiceImpl struct {
	oauthRepo     oauthEnt.OAuthTokenRepository
	oauthProvider goth.Provider
}

func NewOAuthService(repo oauthEnt.OAuthTokenRepository, googleOAuthConf *config.GoogleOAuthConfig) *OAuthServiceImpl {

	scopes := []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"}

	googleProvider := google.New(googleOAuthConf.ClientID, googleOAuthConf.ClientSecret, googleOAuthConf.ClientCallbackURL, scopes...)
	googleProvider.SetAccessType("offline")
	googleProvider.SetPrompt("consent")
	googleProvider.SetName("google")

	goth.UseProviders(googleProvider)

	return &OAuthServiceImpl{oauthRepo: repo, oauthProvider: googleProvider}
}

// InitAuth method save user token data to Redis DB
func (s *OAuthServiceImpl) InitAuth(ctx context.Context, user *goth.User) error {

	existingUser, err := s.oauthRepo.GetToken(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("%s failed to get token: %w", operationName, err)
	}
	if existingUser != nil {
		fmt.Printf("%s user %s already exists, updating...\n", operationName, user.Email)
	}

	// Save or overwrite the user data
	err = s.oauthRepo.SetToken(ctx, user)
	if err != nil {
		return fmt.Errorf("%s failed to save/update token: %w", operationName, err)
	}
	return nil
}

// VerifyAuth method get token data from Redis DB
func (s *OAuthServiceImpl) VerifyAuth(ctx context.Context, identifier string) (*goth.User, error) {
	user, err := s.oauthRepo.GetToken(ctx, identifier)
	if err != nil {
		return nil, fmt.Errorf("%s failed to verify token: %w", operationName, err)
	}
	return user, nil
}

func (s *OAuthServiceImpl) RefreshAuth(ctx context.Context, token string) (string, error) {
	if s.oauthProvider.RefreshTokenAvailable() {
		token, err := s.oauthProvider.RefreshToken(token)
		if err != nil {
			return "", fmt.Errorf("%s failed to refresh token: %w", operationName, err)
		}
		// err = s.oauthRepo.SetToken(ctx, user)
		// if err != nil {
		// 	return "", fmt.Errorf("%s failed to save/update token: %w", operationName, err)
		// }
		return token.AccessToken, nil
	} else {
		return "", fmt.Errorf("%s refresh token not available for this provider", operationName)
	}
}

func (s *OAuthServiceImpl) RevokeAuth(ctx context.Context, c *gin.Context, identifier string) error {
	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		return fmt.Errorf("%s failed to revoke token: %w", operationName, err)
	}
	err = s.oauthRepo.DelToken(ctx, identifier)
	if err != nil {
		return fmt.Errorf("%s failed to delete token: %w", operationName, err)
	}
	return nil
}
