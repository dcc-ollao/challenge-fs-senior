package routes

import (
	"task-management-platform/backend/internal/handlers"
	"task-management-platform/backend/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAPIUserRoutes(r *gin.Engine, h *handlers.APIUserHandler) {
	api := r.Group("/api")
	api.Use(middleware.AuthRequired())

	api.GET("/users", h.ListMinimal)
}
