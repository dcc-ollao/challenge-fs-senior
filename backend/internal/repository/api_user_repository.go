package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"task-management-platform/backend/internal/models"
)

type APIUserRepository interface {
	ListMinimal(ctx context.Context) ([]models.MinimalUser, error)
}

type apiUserRepository struct {
	db *sqlx.DB
}

func NewAPIUserRepository(db *sqlx.DB) APIUserRepository {
	return &apiUserRepository{db: db}
}

func (r *apiUserRepository) ListMinimal(ctx context.Context) ([]models.MinimalUser, error) {
	users := make([]models.MinimalUser, 0)

	query := `
		SELECT id, email
		FROM users
		ORDER BY email ASC
	`
	if err := r.db.SelectContext(ctx, &users, query); err != nil {
		return nil, err
	}
	return users, nil
}
