package server

import (
	"net/http"

	"task-management-platform/backend/internal/handlers"
	"task-management-platform/backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		handlers.RespondOK(c, http.StatusOK, gin.H{"status": "ok"})
	})

	routes.Register(r)

	return r
}
