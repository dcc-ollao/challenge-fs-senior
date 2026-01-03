package routes

import (
	"net/http"

	"task-management-platform/backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		handlers.RespondOK(c, http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
