package auth

import "context"

type Repository[T any] interface {
	// FindByEmail(ctx context.Context, email string) (T, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)
	// Register(ctx context.Context, email, password string) (string, error)
}
