package service

import (
	"context"
	"pethelp-backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type AuthService struct {
	DB        *pgxpool.Pool
	Logger    *zap.Logger
	JwtSecret string
}

func NewAuthService(db *pgxpool.Pool, logger *zap.Logger, jwtSecret string) *AuthService {
	return &AuthService{
		DB:        db,
		Logger:    logger,
		JwtSecret: jwtSecret,
	}
}

// RegisterSpecialist registers a new specialist
func (s *AuthService) RegisterSpecialist(specialist *models.Specialist) error {
	// Implementation here
	return nil
}

func CheckEmailExists(db *pgxpool.Pool, email string) (bool, error) {
	var exists bool
	err := db.QueryRow(context.Background(), 
		"SELECT EXISTS(SELECT 1 FROM specialists WHERE email = $1)", 
		email).Scan(&exists)
	return exists, err
}

func CheckPhoneExists(db *pgxpool.Pool, phone string) (bool, error) {
	var exists bool
	err := db.QueryRow(context.Background(), 
		"SELECT EXISTS(SELECT 1 FROM specialists WHERE phone = $1)", 
		phone).Scan(&exists)
	return exists, err
}
