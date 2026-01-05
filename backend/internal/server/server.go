package server

import (
	"log"
	"net/http"

	"task-management-platform/backend/internal/config"
	"task-management-platform/backend/internal/handlers"
	"task-management-platform/backend/internal/repository"
	"task-management-platform/backend/internal/routes"
	"task-management-platform/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func New(cfg config.Config) *gin.Engine {
	r := gin.Default()

	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)

	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	projectService := services.NewProjectService(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	r.GET("/", func(c *gin.Context) {
		handlers.RespondOK(c, http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/health", func(c *gin.Context) {
		handlers.RespondOK(c, http.StatusOK, gin.H{"status": "ok"})
	})

	routes.Register(r, routes.Dependencies{
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		ProjectHandler: projectHandler,
	})

	return r
}
