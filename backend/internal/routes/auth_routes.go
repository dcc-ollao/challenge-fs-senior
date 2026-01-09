package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"task-management-platform/backend/internal/handlers"
	"task-management-platform/backend/internal/server/middleware"
)

func RegisterAuthRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/change-password", middleware.AuthRequired(), authHandler.ChangePassword)
		auth.GET("/me", middleware.AuthRequired(), func(c *gin.Context) {
			userID, _ := c.Get(middleware.ContextUserIDKey)
			role, _ := c.Get(middleware.ContextRoleKey)

			c.JSON(http.StatusOK, gin.H{
				"userId": userID,
				"role":   role,
			})
		})
		auth.GET(
			"/admin/ping",
			middleware.AuthRequired(),
			middleware.RequireRole("admin"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status": "admin ok",
				})
			},
		)

	}
}
