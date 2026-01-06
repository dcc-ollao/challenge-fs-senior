package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"task-management-platform/backend/internal/models"
	"task-management-platform/backend/internal/services"
)

func writeServiceError(c *gin.Context, err error) {
	switch err {
	case services.ErrBadRequest:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case services.ErrForbidden:
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case services.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}

func mustGetActor(c *gin.Context) models.User {
	v, ok := c.Get("user")
	if !ok {
		return models.User{}
	}
	u, ok := v.(models.User)
	if !ok {
		return models.User{}
	}
	return u
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
