package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
	userID := c.GetString("userId")
	role := c.GetString("role")
	if userID == "" {
		return models.User{}
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return models.User{}
	}

	return models.User{
		ID:   id.String(),
		Role: role,
	}
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
