package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"task-management-platform/backend/internal/handlers/dto"
	"task-management-platform/backend/internal/models"
	"task-management-platform/backend/internal/services"
)

type TaskHandler struct {
	tasks services.TaskService
}

func NewTaskHandler(tasks services.TaskService) *TaskHandler {
	return &TaskHandler{tasks: tasks}
}

func (h *TaskHandler) Create(c *gin.Context) {
	actor := mustGetActor(c)
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var assigneeUUID *uuid.UUID
	if req.AssigneeID != nil {
		parsed, err := uuid.Parse(*req.AssigneeID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignee_id"})
			return
		}
		assigneeUUID = &parsed
	}

	task := &models.Task{
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		AssigneeID:  assigneeUUID,
	}

	if err := h.tasks.Create(c.Request.Context(), actor, task); err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	actor := mustGetActor(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	task, err := h.tasks.GetByID(c.Request.Context(), actor, id)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) ListByProject(c *gin.Context) {
	actor := mustGetActor(c)

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	filters, err := dto.ParseTaskFilters(
		projectID,
		c.Query("status"),
		c.Query("assignee_id"),
		c.Query("limit"),
		c.Query("offset"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query params"})
		return
	}

	tasks, err := h.tasks.List(c.Request.Context(), actor, filters)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) Update(c *gin.Context) {
	actor := mustGetActor(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var assigneeUUID *uuid.UUID
	if req.AssigneeID != nil {
		parsed, err := uuid.Parse(*req.AssigneeID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignee_id"})
			return
		}
		assigneeUUID = &parsed
	}

	task := &models.Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		AssigneeID:  assigneeUUID,
	}

	if err := h.tasks.Update(c.Request.Context(), actor, task); err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	actor := mustGetActor(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.tasks.Delete(c.Request.Context(), actor, id); err != nil {
		writeServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
