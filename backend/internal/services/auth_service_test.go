package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"task-management-platform/backend/internal/models"
)

type fakeUserRepo struct {
	byEmail map[string]*models.User
	byID    map[string]*models.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		byEmail: map[string]*models.User{},
		byID:    map[string]*models.User{},
	}
}

func (r *fakeUserRepo) Create(ctx context.Context, user *models.User) error {
	r.byEmail[user.Email] = user
	r.byID[user.ID] = user
	return nil
}

func (r *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	u, ok := r.byEmail[email]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}

func (r *fakeUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	u, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}

func TestAuthService_Register_Success(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewAuthService(repo)

	user, err := svc.Register(context.Background(), "Test@TEST.com", "1234")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if user.Email != "test@test.com" {
		t.Fatalf("user.Email = %q, want %q", user.Email, "test@test.com")
	}
	if user.Role != "user" {
		t.Fatalf("user.Role = %q, want %q", user.Role, "user")
	}
	if user.PasswordHash == "" {
		t.Fatalf("user.PasswordHash is empty")
	}
	if user.CreatedAt.IsZero() {
		t.Fatalf("user.CreatedAt is zero")
	}
	// sanity: created_at should be near now
	if time.Since(user.CreatedAt) > time.Minute {
		t.Fatalf("user.CreatedAt too old: %v", user.CreatedAt)
	}
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewAuthService(repo)

	_, err := svc.Register(context.Background(), "test@test.com", "1234")
	if err != nil {
		t.Fatalf("first Register() error = %v", err)
	}

	_, err = svc.Register(context.Background(), "test@test.com", "abcd")
	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("second Register() err = %v, want ErrEmailAlreadyExists", err)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewAuthService(repo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	u := &models.User{
		ID:           "u1",
		Email:        "test@test.com",
		PasswordHash: string(hash),
		Role:         "user",
		CreatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(context.Background(), u)

	got, err := svc.Login(context.Background(), "test@test.com", "1234")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if got.ID != "u1" {
		t.Fatalf("got.ID = %q, want %q", got.ID, "u1")
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewAuthService(repo)

	hash, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	u := &models.User{
		ID:           "u1",
		Email:        "test@test.com",
		PasswordHash: string(hash),
		Role:         "user",
		CreatedAt:    time.Now().UTC(),
	}
	_ = repo.Create(context.Background(), u)

	_, err := svc.Login(context.Background(), "test@test.com", "WRONG")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() err = %v, want ErrInvalidCredentials", err)
	}

	_, err = svc.Login(context.Background(), "missing@test.com", "1234")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() err = %v, want ErrInvalidCredentials", err)
	}
}
