package jwt

import (
	"os"
	"testing"
)

func TestGenerateAndParseToken_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_TTL_MINUTES", "60")

	userID := "user-123"
	role := "member"

	token, err := GenerateToken(userID, role)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Fatalf("GenerateToken() returned empty token")
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Fatalf("claims.UserID = %q, want %q", claims.UserID, userID)
	}
	if claims.Role != role {
		t.Fatalf("claims.Role = %q, want %q", claims.Role, role)
	}
	if claims.ExpiresAt == nil {
		t.Fatalf("claims.ExpiresAt is nil")
	}
	if claims.IssuedAt == nil {
		t.Fatalf("claims.IssuedAt is nil")
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_TTL_MINUTES", "60")

	token, err := GenerateToken("user-123", "member")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// breaks the token
	bad := token + "x"

	_, err = ParseToken(bad)
	if err == nil {
		t.Fatalf("ParseToken() expected error, got nil")
	}
	if err != ErrInvalidToken {
		t.Fatalf("ParseToken() error = %v, want %v", err, ErrInvalidToken)
	}
}

func TestGenerateToken_MissingSecret(t *testing.T) {
	_ = os.Unsetenv("JWT_SECRET")

	_, err := GenerateToken("user-123", "member")
	if err == nil {
		t.Fatalf("GenerateToken() expected error, got nil")
	}
}
