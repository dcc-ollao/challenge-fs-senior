package routes

import (
	"github.com/gin-gonic/gin"

	"task-management-platform/backend/internal/handlers"
	"task-management-platform/backend/internal/server/middleware"
)

func RegisterAdminRoutes(r *gin.Engine, exportHandler *handlers.AdminExportHandler) {
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthRequired(), middleware.RequireRole("admin"))
	{
		admin.GET("/export", exportHandler.ExportAll)
	}
}
