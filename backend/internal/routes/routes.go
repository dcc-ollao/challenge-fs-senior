package routes

import (
	"github.com/gin-gonic/gin"

	"task-management-platform/backend/internal/handlers"
)

type Dependencies struct {
	AuthHandler        *handlers.AuthHandler
	UserHandler        *handlers.UserHandler
	ProjectHandler     *handlers.ProjectHandler
	TaskHandler        *handlers.TaskHandler
	APIUserHandler     *handlers.APIUserHandler
	AdminExportHandler *handlers.AdminExportHandler
}

func Register(r *gin.Engine, deps Dependencies) {
	if deps.AuthHandler != nil {
		RegisterAuthRoutes(r, deps.AuthHandler)
	}

	if deps.UserHandler != nil {
		RegisterUserRoutes(r, deps.UserHandler)
	}

	if deps.ProjectHandler != nil {
		RegisterProjectRoutes(r, deps.ProjectHandler)
	}

	if deps.TaskHandler != nil {
		RegisterTaskRoutes(r, deps.TaskHandler)
	}

	if deps.APIUserHandler != nil {
		RegisterAPIUserRoutes(r, deps.APIUserHandler)
	}

	if deps.AdminExportHandler != nil {
		RegisterAdminRoutes(r, deps.AdminExportHandler)
	}
}
