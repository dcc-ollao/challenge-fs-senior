package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"task-management-platform/backend/internal/models"
)

type TaskFilters struct {
	ProjectID  *uuid.UUID
	AssigneeID *uuid.UUID
	Status     *string
	Limit      int
	Offset     int
}

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error)
	List(ctx context.Context, filters TaskFilters) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type taskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	if task.ID == uuid.Nil {
		task.ID = uuid.New()
	}
	query := `
		INSERT INTO tasks (id, project_id, title, description, status, assignee_id)
		VALUES (:id, :project_id, :title, :description, :status, :assignee_id)
	`
	_, err := r.db.NamedExecContext(ctx, query, task)
	return err
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	var task models.Task
	query := `SELECT * FROM tasks WHERE id = $1`
	if err := r.db.GetContext(ctx, &task, query, id); err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) List(ctx context.Context, filters TaskFilters) ([]models.Task, error) {
	var tasks []models.Task

	conditions := []string{"1=1"}
	args := map[string]interface{}{}

	if filters.ProjectID != nil {
		conditions = append(conditions, "project_id = :project_id")
		args["project_id"] = *filters.ProjectID
	}
	if filters.AssigneeID != nil {
		conditions = append(conditions, "assignee_id = :assignee_id")
		args["assignee_id"] = *filters.AssigneeID
	}
	if filters.Status != nil {
		conditions = append(conditions, "status = :status")
		args["status"] = *filters.Status
	}

	query := `
		SELECT *
		FROM tasks
		WHERE ` + strings.Join(conditions, " AND ") + `
		ORDER BY created_at DESC
		LIMIT :limit OFFSET :offset
	`

	args["limit"] = filters.Limit
	args["offset"] = filters.Offset

	rows, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task
		if err := rows.StructScan(&t); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	query := `
		UPDATE tasks
		SET title = :title,
		    description = :description,
		    status = :status,
		    assignee_id = :assignee_id,
		    updated_at = NOW()
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, task)
	return err
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	return err
}
