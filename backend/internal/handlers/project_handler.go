package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"task-management-platform/backend/internal/handlers/dto"
	"task-management-platform/backend/internal/repository"
	"task-management-platform/backend/internal/services"
)

type ProjectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) Register(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	{
		projects.POST("", h.create)
		projects.GET("", h.list)
		projects.GET("/:id", h.getByID)
		projects.PUT("/:id", h.update)
		projects.DELETE("/:id", h.delete)
	}
}

func (h *ProjectHandler) create(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userId")

	p, err := h.service.Create(c.Request.Context(), userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func (h *ProjectHandler) list(c *gin.Context) {
	userID := c.GetString("userId")
	role := c.GetString("role")

	projects, err := h.service.List(c.Request.Context(), userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) getByID(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")
	role := c.GetString("role")

	p, err := h.service.GetByID(c.Request.Context(), userID, role, id)
	if err != nil {
		switch err {
		case services.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case repository.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) update(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")
	role := c.GetString("role")

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateName(c.Request.Context(), userID, role, id, req.Name); err != nil {
		switch err {
		case services.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case repository.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProjectHandler) delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userId")
	role := c.GetString("role")

	if err := h.service.Delete(c.Request.Context(), userID, role, id); err != nil {
		switch err {
		case services.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case repository.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
