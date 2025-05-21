package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/markbates/goth"
	"github.com/redis/go-redis/v9"
)

const (
	userKeyPrefix  = "email:"
	userExpiration = 1 * time.Hour
	operationName  = "oauth_token_repo:"
)

type OAuthTokenRepoImpl struct {
	redisCl *redis.Client
}

func NewOAuthTokenRepo(cli *redis.Client) *OAuthTokenRepoImpl {
	return &OAuthTokenRepoImpl{redisCl: cli}
}

// SetToken method set data to Redis DB
func (r *OAuthTokenRepoImpl) SetToken(ctx context.Context, user *goth.User) error {

	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("%s failed to marshal data: %w", operationName, err)
	}

	key := fmt.Sprintf("%s%s", userKeyPrefix, user.Email)
	err = r.redisCl.Set(ctx, key, userJSON, userExpiration).Err()
	if err != nil {
		return fmt.Errorf("%s failed to save to Redis: %w", operationName, err)
	}
	return nil
}

// GetToken method get data from Redis DB
func (r *OAuthTokenRepoImpl) GetToken(ctx context.Context, identifier string) (*goth.User, error) {
	key := fmt.Sprintf("%s%s", userKeyPrefix, identifier)
	userJSON, err := r.redisCl.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("%s identifier not found in Redis", operationName)
	} else if err != nil {
		return nil, fmt.Errorf("%s failed to get user from Redis: %w", operationName, err)
	}

	var user goth.User
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return nil, fmt.Errorf("%s failed to unmarshal data: %w", operationName, err)
	}
	return &user, nil
}

// DelToken method delete data from Redis DB
func (r *OAuthTokenRepoImpl) DelToken(ctx context.Context, identifier string) error {
	key := fmt.Sprintf("%s%s", userKeyPrefix, identifier)
	err := r.redisCl.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("%s failed to delete data from Redis: %w", operationName, err)
	}
	return nil
}
