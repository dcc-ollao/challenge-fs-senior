package server

import (
	"log"
	"net/http"

	"task-management-platform/backend/internal/config"
	"task-management-platform/backend/internal/handlers"
	"task-management-platform/backend/internal/repository"
	"task-management-platform/backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func New(cfg config.Config) *gin.Engine {
	r := gin.Default()

	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	_ = db

	r.GET("/", func(c *gin.Context) {
		handlers.RespondOK(c, http.StatusOK, gin.H{"status": "ok"})
	})

	routes.Register(r, routes.Dependencies{})

	return r
}
