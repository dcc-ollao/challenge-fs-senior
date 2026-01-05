package dto

type CreateProjectRequest struct {
	Name string `json:"name" binding:"required,min=1,max=120"`
}

type UpdateProjectRequest struct {
	Name string `json:"name" binding:"required,min=1,max=120"`
}
