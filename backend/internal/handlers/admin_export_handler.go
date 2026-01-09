package handlers

import (
	"archive/zip"
	"encoding/csv"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"task-management-platform/backend/internal/repository"
)

type AdminExportHandler struct {
	userRepo    *repository.UserRepository
	projectRepo repository.ProjectRepository
	taskRepo    repository.TaskRepository
}

func NewAdminExportHandler(
	userRepo *repository.UserRepository,
	projectRepo repository.ProjectRepository,
	taskRepo repository.TaskRepository,
) *AdminExportHandler {
	return &AdminExportHandler{
		userRepo:    userRepo,
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
	}
}

func (h *AdminExportHandler) ExportAll(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.userRepo.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to export users"})
		return
	}

	projects, err := h.projectRepo.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to export projects"})
		return
	}

	tasks, err := h.taskRepo.ListAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to export tasks"})
		return
	}

	filename := "export-" + time.Now().UTC().Format("20060102-150405") + ".zip"
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)

	zw := zip.NewWriter(c.Writer)
	defer zw.Close()

	// meta.csv
	{
		f, _ := zw.Create("meta.csv")
		w := csv.NewWriter(f)
		_ = w.Write([]string{"exported_at_utc"})
		_ = w.Write([]string{time.Now().UTC().Format(time.RFC3339)})
		w.Flush()
	}

	// users.csv
	{
		f, _ := zw.Create("users.csv")
		w := csv.NewWriter(f)
		_ = w.Write([]string{"id", "email", "role", "created_at"})
		for _, u := range users {
			_ = w.Write([]string{
				u.ID,
				u.Email,
				u.Role,
				u.CreatedAt.Format(time.RFC3339),
			})
		}
		w.Flush()
	}

	// projects.csv
	{
		f, _ := zw.Create("projects.csv")
		w := csv.NewWriter(f)
		_ = w.Write([]string{"id", "name", "owner_id", "created_at"})
		for _, p := range projects {
			_ = w.Write([]string{
				p.ID,
				p.Name,
				p.OwnerID,
				p.CreatedAt.Format(time.RFC3339),
			})
		}
		w.Flush()
	}

	// tasks.csv
	{
		f, _ := zw.Create("tasks.csv")
		w := csv.NewWriter(f)
		_ = w.Write([]string{"id", "project_id", "title", "description", "status", "assignee_id", "created_at", "updated_at"})
		for _, t := range tasks {
			assignee := ""
			if t.AssigneeID != nil {
				assignee = t.AssigneeID.String()
			}
			_ = w.Write([]string{
				t.ID.String(),
				t.ProjectID.String(),
				t.Title,
				t.Description,
				t.Status,
				assignee,
				t.CreatedAt.Format(time.RFC3339),
				t.UpdatedAt.Format(time.RFC3339),
			})
		}
		w.Flush()
	}
}
