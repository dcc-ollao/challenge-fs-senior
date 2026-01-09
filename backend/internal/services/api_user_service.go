package services

import (
	"context"

	"task-management-platform/backend/internal/models"
	"task-management-platform/backend/internal/repository"
)

type APIUserService interface {
	ListMinimal(ctx context.Context) ([]models.MinimalUser, error)
}

type apiUserService struct {
	repo repository.APIUserRepository
}

func NewAPIUserService(repo repository.APIUserRepository) APIUserService {
	return &apiUserService{repo: repo}
}

func (s *apiUserService) ListMinimal(ctx context.Context) ([]models.MinimalUser, error) {
	return s.repo.ListMinimal(ctx)
}
