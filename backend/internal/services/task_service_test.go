package services

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"

	"task-management-platform/backend/internal/models"
	"task-management-platform/backend/internal/repository"
)

type fakeTaskRepo struct {
	tasks map[uuid.UUID]models.Task
}

func newFakeTaskRepo() *fakeTaskRepo {
	return &fakeTaskRepo{tasks: make(map[uuid.UUID]models.Task)}
}

func (f *fakeTaskRepo) Create(ctx context.Context, task *models.Task) error {
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = task.CreatedAt
	}
	f.tasks[task.ID] = *task
	return nil
}

func (f *fakeTaskRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	t, ok := f.tasks[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return &t, nil
}

func (f *fakeTaskRepo) List(ctx context.Context, filters repository.TaskFilters) ([]models.Task, error) {
	out := make([]models.Task, 0)

	for _, t := range f.tasks {
		if filters.ProjectID != nil && t.ProjectID != *filters.ProjectID {
			continue
		}
		if filters.AssigneeID != nil {
			if t.AssigneeID == nil || *t.AssigneeID != *filters.AssigneeID {
				continue
			}
		}
		if filters.Status != nil && t.Status != *filters.Status {
			continue
		}
		out = append(out, t)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})

	start := filters.Offset
	if start < 0 {
		start = 0
	}
	end := start + filters.Limit
	if filters.Limit <= 0 {
		end = start
	}

	if start >= len(out) {
		return []models.Task{}, nil
	}
	if end > len(out) {
		end = len(out)
	}

	return out[start:end], nil
}

func (f *fakeTaskRepo) Update(ctx context.Context, task *models.Task) error {
	existing, ok := f.tasks[task.ID]
	if ok {
		if task.CreatedAt.IsZero() {
			task.CreatedAt = existing.CreatedAt
		}
		if task.ProjectID == uuid.Nil {
			task.ProjectID = existing.ProjectID
		}
	}
	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = time.Now()
	}

	f.tasks[task.ID] = *task
	return nil
}

func (f *fakeTaskRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(f.tasks, id)
	return nil
}

type fakeProjectRepoForTasks struct{}

func (f *fakeProjectRepoForTasks) Create(ctx context.Context, project *models.Project) error {
	return nil
}
func (f *fakeProjectRepoForTasks) GetByID(ctx context.Context, id string) (*models.Project, error) {
	return &models.Project{}, nil
}
func (f *fakeProjectRepoForTasks) ListByOwner(ctx context.Context, ownerID string) ([]models.Project, error) {
	return []models.Project{}, nil
}
func (f *fakeProjectRepoForTasks) UpdateName(ctx context.Context, id string, name string) error {
	return nil
}
func (f *fakeProjectRepoForTasks) Delete(ctx context.Context, id string) error { return nil }

func ptrUUID(u uuid.UUID) *uuid.UUID { return &u }

func mustCreateTask(t *testing.T, repo *fakeTaskRepo, task models.Task) {
	t.Helper()
	tmp := task
	if tmp.ID == uuid.Nil {
		tmp.ID = uuid.New()
	}
	if err := repo.Create(context.Background(), &tmp); err != nil {
		t.Fatalf("create task: %v", err)
	}
}

func TestTaskService_MemberCannotUpdateUnassignedTask(t *testing.T) {
	taskRepo := newFakeTaskRepo()
	projectRepo := &fakeProjectRepoForTasks{}
	svc := NewTaskService(taskRepo, projectRepo)

	taskID := uuid.New()
	projectID := uuid.New()

	mustCreateTask(t, taskRepo, models.Task{
		ID:        taskID,
		ProjectID: projectID,
		Title:     "orig",
		Status:    "todo",
		CreatedAt: time.Now().Add(-10 * time.Minute),
	})

	member := models.User{
		ID:   uuid.New().String(),
		Role: "user",
	}

	err := svc.Update(context.Background(), member, &models.Task{
		ID:          taskID,
		Title:       "updated",
		Description: "x",
		Status:      "done",
		AssigneeID:  nil,
	})

	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestTaskService_MemberCannotReassignTaskToAnotherUser(t *testing.T) {
	taskRepo := newFakeTaskRepo()
	projectRepo := &fakeProjectRepoForTasks{}
	svc := NewTaskService(taskRepo, projectRepo)

	taskID := uuid.New()
	projectID := uuid.New()
	memberID := uuid.New()
	otherID := uuid.New()

	mustCreateTask(t, taskRepo, models.Task{
		ID:         taskID,
		ProjectID:  projectID,
		Title:      "orig",
		Status:     "todo",
		AssigneeID: ptrUUID(memberID),
		CreatedAt:  time.Now().Add(-10 * time.Minute),
	})

	member := models.User{
		ID:   memberID.String(),
		Role: "user",
	}

	err := svc.Update(context.Background(), member, &models.Task{
		ID:          taskID,
		Title:       "updated",
		Description: "x",
		Status:      "in_progress",
		AssigneeID:  ptrUUID(otherID),
	})

	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestTaskService_AdminCanUpdateAnyTask(t *testing.T) {
	taskRepo := newFakeTaskRepo()
	projectRepo := &fakeProjectRepoForTasks{}
	svc := NewTaskService(taskRepo, projectRepo)

	taskID := uuid.New()
	projectID := uuid.New()

	mustCreateTask(t, taskRepo, models.Task{
		ID:        taskID,
		ProjectID: projectID,
		Title:     "orig",
		Status:    "todo",
		CreatedAt: time.Now().Add(-10 * time.Minute),
	})

	admin := models.User{
		ID:   uuid.New().String(),
		Role: "admin",
	}

	err := svc.Update(context.Background(), admin, &models.Task{
		ID:          taskID,
		Title:       "updated",
		Description: "x",
		Status:      "done",
		AssigneeID:  nil,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskService_ListFiltersByStatus(t *testing.T) {
	taskRepo := newFakeTaskRepo()
	projectRepo := &fakeProjectRepoForTasks{}
	svc := NewTaskService(taskRepo, projectRepo)

	projectID := uuid.New()

	mustCreateTask(t, taskRepo, models.Task{
		ID:        uuid.New(),
		ProjectID: projectID,
		Title:     "t1",
		Status:    "todo",
		CreatedAt: time.Now().Add(-20 * time.Minute),
	})
	mustCreateTask(t, taskRepo, models.Task{
		ID:        uuid.New(),
		ProjectID: projectID,
		Title:     "t2",
		Status:    "done",
		CreatedAt: time.Now().Add(-10 * time.Minute),
	})

	status := "done"
	filters := repository.TaskFilters{
		ProjectID: &projectID,
		Status:    &status,
		Limit:     20,
		Offset:    0,
	}

	admin := models.User{ID: uuid.New().String(), Role: "admin"}

	tasks, err := svc.List(context.Background(), admin, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Status != "done" {
		t.Fatalf("expected status done, got %s", tasks[0].Status)
	}
}

func TestTaskService_ListRespectsPagination(t *testing.T) {
	taskRepo := newFakeTaskRepo()
	projectRepo := &fakeProjectRepoForTasks{}
	svc := NewTaskService(taskRepo, projectRepo)

	projectID := uuid.New()

	now := time.Now()
	mustCreateTask(t, taskRepo, models.Task{
		ID:        uuid.New(),
		ProjectID: projectID,
		Title:     "t1",
		Status:    "todo",
		CreatedAt: now.Add(-30 * time.Minute),
	})
	mustCreateTask(t, taskRepo, models.Task{
		ID:        uuid.New(),
		ProjectID: projectID,
		Title:     "t2",
		Status:    "todo",
		CreatedAt: now.Add(-20 * time.Minute),
	})
	mustCreateTask(t, taskRepo, models.Task{
		ID:        uuid.New(),
		ProjectID: projectID,
		Title:     "t3",
		Status:    "todo",
		CreatedAt: now.Add(-10 * time.Minute),
	})

	filters := repository.TaskFilters{
		ProjectID: &projectID,
		Limit:     1,
		Offset:    1,
	}

	admin := models.User{ID: uuid.New().String(), Role: "admin"}

	tasks, err := svc.List(context.Background(), admin, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Title != "t2" {
		t.Fatalf("expected second newest task title t2, got %s", tasks[0].Title)
	}
}
