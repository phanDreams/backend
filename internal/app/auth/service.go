package appauth

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"pethelp-backend/internal/database/postgres"
	dom "pethelp-backend/internal/domain/auth"
)

type AuthService struct {
	storage *postgres.Storage
	logger *zap.Logger
	hasher          dom.PasswordHasher
    validator       dom.PasswordValidator
	defaultTimeout  time.Duration
	table           string
	jwtSecret       []byte
	cache          *redis.Client
	accessTTL       time.Duration
    refreshTTL      time.Duration
}

type RedisClient = *redis.Client

func NewAuthService(storage *postgres.Storage,
    logger *zap.Logger,
    hasher dom.PasswordHasher,
    validator dom.PasswordValidator,
    jwtSecret string,
    table string,
    cache RedisClient) *AuthService {
	if storage == nil {
		logger.Fatal("Database storage is nil")
	}
	return &AuthService{ 
		storage:        storage,
        logger:         logger,
        hasher:         hasher,
        validator:      validator,
        defaultTimeout: 5 * time.Second,
        table:          table,
        jwtSecret:      []byte(jwtSecret),
        cache:          cache,
        accessTTL:      15 * time.Minute,
        refreshTTL:     7 * 24 * time.Hour,
	}
}

func (s *AuthService) Register(ctx context.Context, model dom.Persistable, pswValidator dom.Credentialed, password string) error {

	if err := s.validator.ValidatePassword(password); err != nil {
        return err
    }

	//hash password
	hashedPassword, err := s.hasher.Hash(password)

	if err != nil {
		s.logger.Error(dom.ErrFailedToHashPassword.Error(), zap.Error(err))
		return err
	}

	pswValidator.SetPasswordHash(hashedPassword)
    
	//build the INSERT, returning ID
	cols := model.Columns()
	placeholders := make([]string, len(cols))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}


	query := fmt.Sprintf(
        "INSERT INTO %s (%s) VALUES (%s) RETURNING id",
        model.TableName(),
        strings.Join(cols, ", "),
        strings.Join(placeholders, ", "),
    )

	//run 
   timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := s.storage.DB().QueryRow(timeoutCtx, query, model.Values()...)
	var id int64
	   if err := row.Scan(&id); err != nil {
	      s.logger.Error("Failed to insert "+model.TableName(), zap.Error(err))
	      return err
	  }
	  model.SetID(id)
	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	//Retrieve stored hash and user ID
	timeoutCtx, cancel := context.WithTimeout(ctx, s.defaultTimeout)
	defer cancel()

	query := fmt.Sprintf(
        "SELECT id, password_hash FROM %s WHERE email = $1",
        s.table,
    )
    var (
        id         int64
        pwHash     string
    )

	row := s.storage.DB().QueryRow(timeoutCtx, query, email)
	if err := row.Scan(&id, &pwHash); err != nil {
        s.logger.Warn("login failed: user not found", zap.String("email", email), zap.Error(err))
        return "", "", dom.ErrInvalidCredentials
    }

	//verify password
	if err := s.hasher.Compare(pwHash, password); err != nil {
		s.logger.Warn("login failed: bad credentials", zap.String("email", email), zap.Error(err))
        return "", "", dom.ErrInvalidCredentials
	}

	subj := strconv.FormatInt(id, 10)

	//generate tokens
	accessToken, err = s.GenerateJWT(subj, s.accessTTL)
	if err != nil {
        s.logger.Error("could not sign access token", zap.Error(err))
        return "", "", dom.ErrInvalidToken
    }
	refreshToken, err = s.GenerateJWT(subj, s.refreshTTL)
    if err != nil {
        s.logger.Error("could not sign refresh token", zap.Error(err))
        return "", "", dom.ErrInvalidToken
    }


	return accessToken, refreshToken, nil
}
