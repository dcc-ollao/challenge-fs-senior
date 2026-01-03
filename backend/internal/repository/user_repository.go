package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"task-management-platform/backend/internal/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, created_at)
		VALUES (:id, :email, :password_hash, :role, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
