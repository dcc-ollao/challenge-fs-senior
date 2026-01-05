package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"task-management-platform/backend/internal/models"
)

type ProjectRepository interface {
	Create(context.Context, *models.Project) error
	GetByID(context.Context, string) (*models.Project, error)
	ListByOwner(context.Context, string) ([]models.Project, error)
	UpdateName(context.Context, string, string) error
	Delete(context.Context, string) error
}

type ProjectService struct {
	repo ProjectRepository
}

func NewProjectService(repo ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ctx context.Context, ownerID string, name string) (*models.Project, error) {
	p := &models.Project{
		ID:        uuid.NewString(),
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) GetByID(ctx context.Context, requesterID string, requesterRole string, projectID string) (*models.Project, error) {
	p, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if !isAdmin(requesterRole) && p.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	return p, nil
}

func (s *ProjectService) List(ctx context.Context, requesterID string, requesterRole string) ([]models.Project, error) {
	_ = requesterRole
	return s.repo.ListByOwner(ctx, requesterID)
}

func (s *ProjectService) UpdateName(ctx context.Context, requesterID string, requesterRole string, projectID string, name string) error {
	p, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	if !isAdmin(requesterRole) && p.OwnerID != requesterID {
		return ErrForbidden
	}

	return s.repo.UpdateName(ctx, projectID, name)
}

func (s *ProjectService) Delete(ctx context.Context, requesterID string, requesterRole string, projectID string) error {
	p, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	if !isAdmin(requesterRole) && p.OwnerID != requesterID {
		return ErrForbidden
	}

	return s.repo.Delete(ctx, projectID)
}

func isAdmin(role string) bool {
	return role == "admin"
}
