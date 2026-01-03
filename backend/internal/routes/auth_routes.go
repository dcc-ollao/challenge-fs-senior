package routes

import (
	"github.com/gin-gonic/gin"

	"task-management-platform/backend/internal/handlers"
)

func RegisterAuthRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}
