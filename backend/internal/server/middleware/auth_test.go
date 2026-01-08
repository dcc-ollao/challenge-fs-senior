package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	jwtutil "task-management-platform/backend/pkg/jwt"
)

func TestAuthRequired_MissingAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/protected", AuthRequired(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthRequired_InvalidAuthorizationHeaderFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/protected", AuthRequired(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token abcdef")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_TTL_MINUTES", "60")

	r := gin.New()
	r.GET("/protected", AuthRequired(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthRequired_ValidToken_SetsContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_TTL_MINUTES", "60")

	token, err := jwtutil.GenerateToken("user-123", "member")
	if err != nil {
		t.Fatalf("GenerateToken error = %v", err)
	}

	r := gin.New()
	r.GET("/protected", AuthRequired(), func(c *gin.Context) {
		userID, _ := c.Get(ContextUserIDKey)
		role, _ := c.Get(ContextRoleKey)

		c.JSON(http.StatusOK, gin.H{
			"userId": userID,
			"role":   role,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if body["userId"] != "user-123" {
		t.Fatalf("userId = %v, want %v", body["userId"], "user-123")
	}
	if body["role"] != "user" {
		t.Fatalf("role = %v, want %v", body["role"], "user")
	}
}
