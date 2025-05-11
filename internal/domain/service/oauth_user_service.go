package service

import (
	"context"
	"fmt"

	"pethelp-backend/internal/config"
	"pethelp-backend/internal/repository"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

type OAuthUserService struct {
	userRepo      *repository.OAuthUserRepository
	oauthProvider goth.Provider
}

func NewOAuthService(userRepo *repository.OAuthUserRepository, googleOAuthConf config.GoogleOAuthConfig) *OAuthUserService {

	googleProvider := google.New(googleOAuthConf.ClientID, googleOAuthConf.ClientSecret, googleOAuthConf.ClientCallbackURL)
	googleProvider.SetAccessType("offline")
	googleProvider.SetPrompt("consent")
	googleProvider.SetName("google")

	goth.UseProviders(googleProvider)

	return &OAuthUserService{userRepo: userRepo, oauthProvider: googleProvider}
}

func (s *OAuthUserService) SetOrUpdateUser(ctx context.Context, user *goth.User) error {
	// Check if the user already exists (e.g., by ID)
	existingUser, err := s.userRepo.GetUser(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %w", err)
	}

	if existingUser != nil {
		fmt.Printf("User with ID %s already exists, updating...\n", user.Email)
	}

	// Save or overwrite the user data
	err = s.userRepo.SetUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to save/update user: %w", err)
	}
	return nil
}

func (s *OAuthUserService) GetUserByEmail(ctx context.Context, userEmail string) (*goth.User, error) {
	user, err := s.userRepo.GetUser(ctx, userEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}
