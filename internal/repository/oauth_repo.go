package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/markbates/goth"
	"github.com/redis/go-redis/v9"
)

const userKeyPrefix = "email:"
const userExpiration = 24 * time.Hour // Example expiration time
const operationName = "oauth_repo"

type OAuthUserRepository struct {
	redisClient *redis.Client
}

func NewOAuthUserRepository(redisClient *redis.Client) *OAuthUserRepository {
	return &OAuthUserRepository{redisClient: redisClient}
}

func (r *OAuthUserRepository) SetUser(ctx context.Context, user *goth.User) error {

	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("%s failed to marshal user data: %w", operationName, err)
	}

	key := fmt.Sprintf("%s%s", userKeyPrefix, user.Email) // Use a unique identifier like user.ID
	err = r.redisClient.Set(ctx, key, userJSON, userExpiration).Err()
	if err != nil {
		return fmt.Errorf("%s failed to save user to Redis: %w", operationName, err)
	}
	return nil
}

func (r *OAuthUserRepository) GetUser(ctx context.Context, userEmail string) (*goth.User, error) {
	key := fmt.Sprintf("%s%s", userKeyPrefix, userEmail)
	userJSON, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // User not found
	} else if err != nil {
		return nil, fmt.Errorf("%s failed to get user from Redis: %w", operationName, err)
	}

	var user goth.User
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return nil, fmt.Errorf("%s failed to unmarshal user data: %w", operationName, err)
	}
	return &user, nil
}
