package server

import (
	"log"
	"net/http"
	"os"
	"strings"
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
	origins := []string{
		"http://localhost:5173",
		"http://localhost:3000",
	}

	if v := os.Getenv("CORS_ORIGINS"); v != "" {
		origins = strings.Split(v, ",")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Disposition"}, // 👈 importante para export ZIP/CSV
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
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
