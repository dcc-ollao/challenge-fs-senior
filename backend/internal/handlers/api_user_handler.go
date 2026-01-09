package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"task-management-platform/backend/internal/services"
)

type APIUserHandler struct {
	svc services.APIUserService
}

func NewAPIUserHandler(svc services.APIUserService) *APIUserHandler {
	return &APIUserHandler{svc: svc}
}

func (h *APIUserHandler) ListMinimal(c *gin.Context) {
	users, err := h.svc.ListMinimal(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, users)
}
