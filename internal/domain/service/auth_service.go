package service

import (
	"context"
	"fmt"
	"time"

	"pethelp-backend/internal/domain/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// handle auth-related operations
type AuthService struct {
	DB        *pgxpool.Pool
	Logger    *zap.Logger
	JwtSecret string
}

// create a new instance os AuthService
func NewAuthService(db *pgxpool.Pool, logger *zap.Logger, jwtSecret string) (*AuthService, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is empty")
	}

	return &AuthService{
		DB:        db,
		Logger:    logger,
		JwtSecret: jwtSecret,
	}, nil
}

// register a new specialist
func (s *AuthService) RegisterSpecialist(specialist *models.Specialist) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(specialist.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("Failed to hash password", zap.Error(err))
		return err
	}

	// store the hashed password in the database
	passwordHash := string(hashedPassword)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlStmt := `
		INSERT INTO specialists (name, family_name, phone, email, password_hash) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err = s.DB.QueryRow(ctx, sqlStmt,
		specialist.Name,
		specialist.FamilyName,
		specialist.Phone,
		specialist.Email,
		passwordHash).Scan(&specialist.ID)

	if err != nil {
		s.Logger.Error("Failed to insert specialist", zap.Error(err))
		return err
	}
	return nil
}

func (s *AuthService) CheckEmailExists(email string) (bool, error) {
	const operationName = "CheckEmailExists" //  For consistent logging

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //  Good default.  Consider making this a config option.
	defer cancel()

	const query = `SELECT EXISTS(SELECT 1 FROM specialists WHERE email = $1 AND is_deleted = false)`

	var exists bool
	err := s.DB.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Email does not exist (which is a valid case, not an error).
			return false, nil //  Return false, nil for non-existent email
		}
		// Wrap the error with context.
		err = fmt.Errorf("%s: failed to check email existence: %w", operationName, err)
		s.Logger.Error("Database query failed", zap.Error(err), zap.String("email", email))
		return false, err
	}

	s.Logger.Debug("Email existence check successful", zap.String("email", email))
	return exists, nil
}

func (s *AuthService) CheckPhoneExists(phone string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	// Use simple query instead of prepared statement
	err := s.DB.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM specialists WHERE phone = $1 AND is_deleted = false)",
		phone).Scan(&exists)
	if err != nil {
		s.Logger.Error("Failed to check phone existence", zap.Error(err))
		return false, err
	}
	return exists, nil
}

// generate token
func (s *AuthService) GenerateToken(specialist *models.Specialist) (string, error) {
	claims := jwt.MapClaims{
		"specialist_id": specialist.ID,
		"email":         specialist.Email,
		"exp":           time.Now().Add(72 * time.Hour).Unix(), // 72-hour expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JwtSecret))
}
