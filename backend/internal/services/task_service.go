package services

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"task-management-platform/backend/internal/models"
	"task-management-platform/backend/internal/repository"
)

const (
	roleAdmin  = "admin"
	roleMember = "user"
)

type TaskService interface {
	Create(ctx context.Context, actor models.User, task *models.Task) error
	GetByID(ctx context.Context, actor models.User, id uuid.UUID) (*models.Task, error)
	List(ctx context.Context, actor models.User, filters repository.TaskFilters) ([]models.Task, error)
	Update(ctx context.Context, actor models.User, task *models.Task) error
	Delete(ctx context.Context, actor models.User, id uuid.UUID) error
}

type taskService struct {
	tasks    repository.TaskRepository
	projects repository.ProjectRepository
}

func NewTaskService(tasks repository.TaskRepository, projects repository.ProjectRepository) TaskService {
	return &taskService{
		tasks:    tasks,
		projects: projects,
	}
}

func (s *taskService) Create(ctx context.Context, actor models.User, task *models.Task) error {
	if err := validateTaskCreate(task); err != nil {
		return err
	}

	if _, err := s.projects.GetByID(ctx, task.ProjectID.String()); err != nil {
		return ErrNotFound
	}

	if actor.Role == roleAdmin {
		return s.tasks.Create(ctx, task)
	}

	if actor.Role == roleMember {
		if task.AssigneeID != nil && task.AssigneeID.String() != actor.ID {
			return ErrForbidden
		}
		return s.tasks.Create(ctx, task)
	}

	return ErrForbidden
}

func (s *taskService) GetByID(ctx context.Context, actor models.User, id uuid.UUID) (*models.Task, error) {
	task, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}

	if actor.Role == roleAdmin {
		return task, nil
	}

	if actor.Role == roleMember {
		return task, nil
	}

	return nil, ErrForbidden
}

func (s *taskService) List(ctx context.Context, actor models.User, filters repository.TaskFilters) ([]models.Task, error) {
	normalizeFilters(&filters)

	if actor.Role == roleAdmin {
		return s.tasks.List(ctx, filters)
	}

	if actor.Role == roleMember {
		return s.tasks.List(ctx, filters)
	}

	return nil, ErrForbidden
}

func (s *taskService) Update(ctx context.Context, actor models.User, task *models.Task) error {
	if err := validateTaskUpdate(task); err != nil {
		return err
	}

	existing, err := s.tasks.GetByID(ctx, task.ID)
	if err != nil {
		return ErrNotFound
	}

	if actor.Role == roleAdmin {
		task.ProjectID = existing.ProjectID
		return s.tasks.Update(ctx, task)
	}

	if actor.Role == roleMember {
		if existing.AssigneeID == nil || existing.AssigneeID.String() != actor.ID {
			return ErrForbidden
		}
		if task.AssigneeID != nil && task.AssigneeID.String() != actor.ID {
			return ErrForbidden
		}

		task.ProjectID = existing.ProjectID
		return s.tasks.Update(ctx, task)
	}

	return ErrForbidden
}

func (s *taskService) Delete(ctx context.Context, actor models.User, id uuid.UUID) error {
	existing, err := s.tasks.GetByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	if actor.Role == roleAdmin {
		return s.tasks.Delete(ctx, existing.ID)
	}

	if actor.Role == roleMember {
		if existing.AssigneeID == nil || existing.AssigneeID.String() != actor.ID {
			return ErrForbidden
		}
		return s.tasks.Delete(ctx, existing.ID)
	}

	return ErrForbidden
}

func validateTaskCreate(task *models.Task) error {
	if task == nil {
		return ErrBadRequest
	}
	task.Title = strings.TrimSpace(task.Title)
	task.Status = strings.TrimSpace(task.Status)

	if task.ProjectID == uuid.Nil {
		return ErrBadRequest
	}
	if task.Title == "" {
		return ErrBadRequest
	}
	if task.Status == "" {
		task.Status = "todo"
	}
	if !isValidStatus(task.Status) {
		return ErrBadRequest
	}
	return nil
}

func validateTaskUpdate(task *models.Task) error {
	if task == nil {
		return ErrBadRequest
	}
	if task.ID == uuid.Nil {
		return ErrBadRequest
	}
	task.Title = strings.TrimSpace(task.Title)
	task.Status = strings.TrimSpace(task.Status)

	if task.Title == "" {
		return ErrBadRequest
	}
	if task.Status == "" || !isValidStatus(task.Status) {
		return ErrBadRequest
	}
	return nil
}

func isValidStatus(s string) bool {
	switch s {
	case "todo", "in_progress", "done":
		return true
	default:
		return false
	}
}

func normalizeFilters(f *repository.TaskFilters) {
	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
}
