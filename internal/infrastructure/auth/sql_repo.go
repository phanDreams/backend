package authinfrastructure

import (
	"context"
	"fmt"
	"pethelp-backend/internal/database/postgres"
	dom "pethelp-backend/internal/domain/auth"

	"go.uber.org/zap"
)

type SQLRepository[E dom.Registrable] struct {
	storage *postgres.Storage
	logger  *zap.Logger
	table   string
}

func NewSQLRepository[E dom.Registrable](storage *postgres.Storage, logger *zap.Logger, table string) dom.Repository[E] {
	return &SQLRepository[E]{
		storage: storage,
		logger:  logger,
		table:   table,
	}
}

// func (r *SQLRepository[E]) FindByEmail(ctx context.Context, email string) (E, error) {
// 	var entity E
// 	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1 AND is_deleted = false", r.table)
// 	err := r.storage.DB().QueryRow(ctx, query, email).Scan(entity.Values()...)
// 	return entity, err
// }

func (r *SQLRepository[E]) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE email = $1 AND is_deleted = false)", r.table)
	if err := r.storage.DB().QueryRow(ctx, query, email).Scan(&exists); err != nil {
		r.logger.Error("email exists check failed", zap.Error(err))
		return false, fmt.Errorf("CheckEmailExists %q: %w", email, err)
	}
	return exists, nil
}

func (r *SQLRepository[E]) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE phone = $1 AND is_deleted = false)", r.table)
	if err := r.storage.DB().QueryRow(ctx, query, phone).Scan(&exists); err != nil {
		r.logger.Error("phone exists check failed", zap.Error(err))
		return false, fmt.Errorf("CheckPhoneExists %q: %w", phone, err)
	}
	return exists, nil
}

// func (r *SQLRepository[E]) Register(ctx context.Context, email, password string) (string, error) {
// 	// This method is not used in the current implementation
// 	return "", nil
// }

