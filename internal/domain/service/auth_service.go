package service

import (
	"context"
	"time"

	"pethelp-backend/internal/domain/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

//handle auth-related operations
type AuthService struct {
	DB        *pgxpool.Pool
	Logger    *zap.Logger
	JwtSecret string
}

//create a new instance os AuthService
func NewAuthService(db *pgxpool.Pool, logger *zap.Logger, jwtSecret string) *AuthService {
	return &AuthService{
		DB:        db,
		Logger:    logger,
		JwtSecret: jwtSecret,
	}
}

//register a new specialist
func (s *AuthService) RegisterSpecialist(specialist *models.Specialist) error {
   hashedPassword, err := bcrypt.GenerateFromPassword([]byte(specialist.Password), bcrypt.DefaultCost)

   if err != nil {
	s.Logger.Error("Failed to hash password", zap.Error(err))
	return err
	}

	// store the hashed password in the database
	specialist.Password = string(hashedPassword)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlStmt := `INSERT INTO specialists (name, family_name, phone, email, password_hash, is_banned, is_deleted, is_active, is_verified, created_at) 
	VALUES ($1, $2, $3, $4, $5, false, false, true, false, NOW())
	RETURNING id`

	err = s.DB.QueryRow(ctx, sqlStmt,specialist.Name, specialist.FamilyName, specialist.Phone, specialist.Email, specialist.Password,
		specialist.Is_banned, specialist.Is_deleted, specialist.Is_active, specialist.Is_verified).Scan(&specialist.ID)

	if err != nil {
		s.Logger.Error("Failed to insert specialist", zap.Error(err))
		return err
	}
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


//generate token
func (s *AuthService) GenerateToken(specialist *models.Specialist) (string, error) {
	claims := jwt.MapClaims{
		"specialist_id": specialist.ID,
		"email":   specialist.Email,
		"exp":     time.Now().Add(72 * time.Hour).Unix(), // 72-hour expiration
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JwtSecret))
}