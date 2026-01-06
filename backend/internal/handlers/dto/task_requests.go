package dto

type CreateTaskRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	AssigneeID  *string `json:"assignee_id"`
}

type UpdateTaskRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Status      string  `json:"status" binding:"required"`
	AssigneeID  *string `json:"assignee_id"`
}
