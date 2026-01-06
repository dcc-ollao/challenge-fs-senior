package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"task-management-platform/backend/internal/models"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id string) (*models.Project, error)
	ListByOwner(ctx context.Context, ownerID string) ([]models.Project, error)
	UpdateName(ctx context.Context, id string, name string) error
	Delete(ctx context.Context, id string) error
}

type projectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *models.Project) error {
	query := `
		INSERT INTO projects (id, name, owner_id, created_at)
		VALUES (:id, :name, :owner_id, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, project)
	return err
}

func (r *projectRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
	var p models.Project

	query := `
		SELECT id, name, owner_id, created_at
		FROM projects
		WHERE id = $1
	`
	if err := r.db.GetContext(ctx, &p, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *projectRepository) ListByOwner(ctx context.Context, ownerID string) ([]models.Project, error) {
	projects := make([]models.Project, 0)

	query := `
		SELECT id, name, owner_id, created_at
		FROM projects
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`
	if err := r.db.SelectContext(ctx, &projects, query, ownerID); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) UpdateName(ctx context.Context, id string, name string) error {
	query := `
		UPDATE projects
		SET name = $2
		WHERE id = $1
	`
	res, err := r.db.ExecContext(ctx, query, id, name)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *projectRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM projects WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}
