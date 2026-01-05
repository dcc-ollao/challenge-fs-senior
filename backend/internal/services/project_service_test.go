package services

import (
	"context"
	"testing"
	"time"

	"task-management-platform/backend/internal/models"
	"task-management-platform/backend/internal/repository"
)

type fakeProjectRepo struct {
	projects map[string]models.Project
}

func newFakeProjectRepo() *fakeProjectRepo {
	return &fakeProjectRepo{
		projects: make(map[string]models.Project),
	}
}

func (f *fakeProjectRepo) Create(ctx context.Context, p *models.Project) error {
	f.projects[p.ID] = *p
	return nil
}

func (f *fakeProjectRepo) GetByID(ctx context.Context, id string) (*models.Project, error) {
	p, ok := f.projects[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return &p, nil
}

func (f *fakeProjectRepo) ListByOwner(ctx context.Context, ownerID string) ([]models.Project, error) {
	out := make([]models.Project, 0)
	for _, p := range f.projects {
		if p.OwnerID == ownerID {
			out = append(out, p)
		}
	}
	return out, nil
}

func (f *fakeProjectRepo) UpdateName(ctx context.Context, id string, name string) error {
	p, ok := f.projects[id]
	if !ok {
		return repository.ErrNotFound
	}
	p.Name = name
	f.projects[id] = p
	return nil
}

func (f *fakeProjectRepo) Delete(ctx context.Context, id string) error {
	if _, ok := f.projects[id]; !ok {
		return repository.ErrNotFound
	}
	delete(f.projects, id)
	return nil
}

func TestProjectOwnerCanAccess(t *testing.T) {
	repo := newFakeProjectRepo()
	svc := NewProjectService(repo)

	p := models.Project{
		ID:        "p1",
		Name:      "project",
		OwnerID:   "user1",
		CreatedAt: time.Now(),
	}
	_ = repo.Create(context.Background(), &p)

	_, err := svc.GetByID(context.Background(), "user1", "member", "p1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProjectNonOwnerForbidden(t *testing.T) {
	repo := newFakeProjectRepo()
	svc := NewProjectService(repo)

	p := models.Project{
		ID:        "p1",
		Name:      "project",
		OwnerID:   "user1",
		CreatedAt: time.Now(),
	}
	_ = repo.Create(context.Background(), &p)

	_, err := svc.GetByID(context.Background(), "user2", "member", "p1")
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAdminCanAccessAnyProject(t *testing.T) {
	repo := newFakeProjectRepo()
	svc := NewProjectService(repo)

	p := models.Project{
		ID:        "p1",
		Name:      "project",
		OwnerID:   "user1",
		CreatedAt: time.Now(),
	}
	_ = repo.Create(context.Background(), &p)

	_, err := svc.GetByID(context.Background(), "admin1", "admin", "p1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateForbidden(t *testing.T) {
	repo := newFakeProjectRepo()
	svc := NewProjectService(repo)

	p := models.Project{
		ID:        "p1",
		Name:      "old",
		OwnerID:   "user1",
		CreatedAt: time.Now(),
	}
	_ = repo.Create(context.Background(), &p)

	err := svc.UpdateName(context.Background(), "user2", "member", "p1", "new")
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestDeleteByOwner(t *testing.T) {
	repo := newFakeProjectRepo()
	svc := NewProjectService(repo)

	p := models.Project{
		ID:        "p1",
		Name:      "project",
		OwnerID:   "user1",
		CreatedAt: time.Now(),
	}
	_ = repo.Create(context.Background(), &p)

	err := svc.Delete(context.Background(), "user1", "member", "p1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	repo := newFakeProjectRepo()
	svc := NewProjectService(repo)

	err := svc.Delete(context.Background(), "user1", "member", "missing")
	if err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	repo := newFakeProjectRepo()
	svc := NewProjectService(repo)

	_, err := svc.GetByID(context.Background(), "user1", "member", "missing")
	if err != repository.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
