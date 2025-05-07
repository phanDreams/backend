package appauth

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	table string
}

func NewAuthService(storage *postgres.Storage, logger *zap.Logger, hasher dom.PasswordHasher, validator dom.PasswordValidator, jwtSecret string, table string) *AuthService {
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
	err = row.Scan(model.SetID)
	if err != nil {
		s.logger.Error("Failed to insert "+model.TableName(), zap.Error(err))
		return err
	}
	return nil
}

// func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
// 	// Implement login logic here
// 	return "", nil
// }
