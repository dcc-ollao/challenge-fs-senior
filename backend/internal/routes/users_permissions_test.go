package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"task-management-platform/backend/internal/server/middleware"
	jwtutil "task-management-platform/backend/pkg/jwt"
)

func TestUsersRoutes_UnauthorizedWhenMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_TTL_MINUTES", "60")

	r := gin.New()

	users := r.Group("/users")
	users.Use(middleware.AuthRequired())
	users.Use(middleware.RequireRole("admin"))
	users.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d. body=%s", w.Code, http.StatusUnauthorized, w.Body.String())
	}
}

func TestUsersRoutes_ForbidsNonAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_TTL_MINUTES", "60")

	memberToken, err := jwtutil.GenerateToken("user-123", "member")
	if err != nil {
		t.Fatalf("GenerateToken error = %v", err)
	}

	r := gin.New()

	users := r.Group("/users")
	users.Use(middleware.AuthRequired())
	users.Use(middleware.RequireRole("admin"))
	users.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Bearer "+memberToken)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d. body=%s", w.Code, http.StatusForbidden, w.Body.String())
	}
}

func TestUsersRoutes_AllowsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_TTL_MINUTES", "60")

	adminToken, err := jwtutil.GenerateToken("admin-1", "admin")
	if err != nil {
		t.Fatalf("GenerateToken error = %v", err)
	}

	r := gin.New()

	users := r.Group("/users")
	users.Use(middleware.AuthRequired())
	users.Use(middleware.RequireRole("admin"))
	users.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body=%s", w.Code, http.StatusOK, w.Body.String())
	}
}
