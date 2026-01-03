package services

import (
	"context"
	"testing"
)

func TestUserService_Delete_CannotDeleteSelf(t *testing.T) {
	svc := &UserService{repo: nil}

	err := svc.Delete(context.Background(), "same-id", "same-id")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
