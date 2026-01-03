package routes

import (
	"task-management-platform/backend/internal/handlers"
	"task-management-platform/backend/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, h *handlers.UserHandler) {
	users := r.Group("/users")
	users.Use(middleware.AuthRequired())
	users.Use(middleware.RequireRole("admin"))

	users.GET("", h.List)
	users.GET("/:id", h.GetByID)
	users.POST("", h.Create)
	users.PUT("/:id", h.Update)
	users.DELETE("/:id", h.Delete)
}
