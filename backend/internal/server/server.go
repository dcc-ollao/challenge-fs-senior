package server

import (
	"log"
	"net/http"
	"time"

	"task-management-platform/backend/internal/config"
	"task-management-platform/backend/internal/handlers"
	"task-management-platform/backend/internal/repository"
	"task-management-platform/backend/internal/routes"
	"task-management-platform/backend/internal/server/middleware"
	"task-management-platform/backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New(cfg config.Config) *gin.Engine {
	r := gin.Default()
	rl := middleware.NewRateLimiter(120, time.Minute)
	r.Use(middleware.RateLimit(rl))
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	apiUserRepo := repository.NewAPIUserRepository(db)

	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	projectService := services.NewProjectService(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	taskService := services.NewTaskService(taskRepo, projectRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	apiUserSvc := services.NewAPIUserService(apiUserRepo)
	apiUserHandler := handlers.NewAPIUserHandler(apiUserSvc)
	adminExportHandler := handlers.NewAdminExportHandler(userRepo, projectRepo, taskRepo)

	r.GET("/", func(c *gin.Context) {
		handlers.RespondOK(c, http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/health", func(c *gin.Context) {
		handlers.RespondOK(c, http.StatusOK, gin.H{"status": "ok"})
	})

	routes.Register(r, routes.Dependencies{
		AuthHandler:        authHandler,
		UserHandler:        userHandler,
		ProjectHandler:     projectHandler,
		TaskHandler:        taskHandler,
		APIUserHandler:     apiUserHandler,
		AdminExportHandler: adminExportHandler,
	})

	return r
}
