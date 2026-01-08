package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	jwtutil "task-management-platform/backend/pkg/jwt"
)

func TestRequireRole_ForbidsNonAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_TTL_MINUTES", "60")

	memberToken, err := jwtutil.GenerateToken("user-123", "user")
	if err != nil {
		t.Fatalf("GenerateToken error = %v", err)
	}

	r := gin.New()
	r.GET("/admin",
		AuthRequired(),
		RequireRole("admin"),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+memberToken)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d. body=%s", w.Code, http.StatusForbidden, w.Body.String())
	}
}

func TestRequireRole_AllowsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_TTL_MINUTES", "60")

	adminToken, err := jwtutil.GenerateToken("admin-1", "admin")
	if err != nil {
		t.Fatalf("GenerateToken error = %v", err)
	}

	r := gin.New()
	r.GET("/admin",
		AuthRequired(),
		RequireRole("admin"),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body=%s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestRequireRole_UnauthorizedWhenMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_TTL_MINUTES", "60")

	r := gin.New()
	r.GET("/admin",
		AuthRequired(),
		RequireRole("admin"),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ok": true})
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
