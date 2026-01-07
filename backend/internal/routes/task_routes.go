package routes

import (
	"task-management-platform/backend/internal/handlers"
	"task-management-platform/backend/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(r *gin.Engine, h *handlers.TaskHandler) {
	api := r.Group("/api")
	api.Use(middleware.AuthRequired())

	api.POST("/projects/:id/tasks", h.Create)
	api.GET("/projects/:id/tasks", h.ListByProject)

	api.GET("/tasks/:id", h.GetByID)
	api.PUT("/tasks/:id", h.Update)
	api.DELETE("/tasks/:id", h.Delete)
}
