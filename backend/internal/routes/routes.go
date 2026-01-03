package routes

import (
	"github.com/gin-gonic/gin"

	"task-management-platform/backend/internal/handlers"
)

type Dependencies struct {
	AuthHandler *handlers.AuthHandler
}

func Register(r *gin.Engine, deps Dependencies) {
	if deps.AuthHandler != nil {
		RegisterAuthRoutes(r, deps.AuthHandler)
	}
}
