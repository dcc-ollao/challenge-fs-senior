package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"task-management-platform/backend/internal/handlers/dto"
	"task-management-platform/backend/internal/models"
	"task-management-platform/backend/internal/repository"
	"task-management-platform/backend/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.service.List(c.Request.Context())
	if err != nil {
		RespondError(c, http.StatusInternalServerError, CodeInternal, "internal error", nil)

		return
	}

	resp := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, dto.FromUser(u))
	}

	RespondOK(c, http.StatusOK, resp)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrNotFound {
			RespondError(c, http.StatusNotFound, CodeNotFound, "user not found", gin.H{"id": id})
			return
		}
		RespondError(c, http.StatusInternalServerError, CodeInternal, "internal error", nil)
		return
	}

	RespondOK(c, http.StatusOK, dto.FromUser(*user))
}

func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(
			c,
			http.StatusBadRequest,
			CodeBadRequest,
			"invalid payload",
			gin.H{"error": err.Error()},
		)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondError(
			c,
			http.StatusInternalServerError,
			CodeInternal,
			"password hashing failed",
			nil,
		)
		return
	}

	now := time.Now()

	user := &models.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		Role:         req.Role,
		PasswordHash: string(hash),
		CreatedAt:    now,
	}

	if err := h.service.Create(c.Request.Context(), user); err != nil {
		RespondError(
			c,
			http.StatusInternalServerError,
			CodeInternal,
			"could not create user",
			nil,
		)
		return
	}

	RespondOK(c, http.StatusCreated, dto.FromUser(*user))
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	actorID := c.GetString("userId")
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Role  string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, CodeBadRequest, "invalid payload", gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		ID:    id,
		Email: req.Email,
		Role:  req.Role,
	}

	if err := h.service.Update(c.Request.Context(), actorID, user); err != nil {
		switch err {
		case repository.ErrNotFound:
			RespondError(c, http.StatusNotFound, CodeNotFound, "user not found", gin.H{"id": id})
			return
		case services.ErrCannotUpdateOwnRole:
			RespondError(c, http.StatusForbidden, CodeForbidden, "cannot update own role", nil)
			return
		default:
			RespondError(c, http.StatusInternalServerError, CodeInternal, "could not update user", nil)
			return
		}
	}

	RespondOK(c, http.StatusOK, gin.H{"status": "updated"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	targetID := c.Param("id")
	actorID := c.GetString("userId")

	if err := h.service.Delete(c.Request.Context(), actorID, targetID); err != nil {
		if err == repository.ErrNotFound {
			RespondError(c, http.StatusNotFound, CodeNotFound, "user not found", gin.H{"id": targetID})
			return
		}
		if err.Error() == "cannot delete own user" {
			RespondError(c, http.StatusBadRequest, CodeBusinessRule, err.Error(), gin.H{"id": targetID})
			return
		}
		RespondError(c, http.StatusInternalServerError, CodeInternal, "could not delete user", nil)
		return
	}

	RespondOK(c, http.StatusOK, gin.H{"status": "deleted"})
}
