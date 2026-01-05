package routes

import (
	"task-management-platform/backend/internal/handlers"
	"task-management-platform/backend/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterProjectRoutes(r *gin.Engine, h *handlers.ProjectHandler) {
	api := r.Group("/api")
	api.Use(middleware.AuthRequired())

	h.Register(api)
}
