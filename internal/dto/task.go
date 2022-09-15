package dto

import (
	"github.com/viniosilva/swordhealth-api/internal/model"
)

type TaskDto struct {
	ID        int              `json:"id" example:"1"`
	CreatedAt string           `json:"created_at,omitempty" example:"1992-08-21 12:03:43"`
	UpdatedAt string           `json:"updated_at,omitempty" example:"1992-08-21 12:03:43"`
	User      UserDto          `json:"user,omitempty"`
	Summary   string           `json:"summary,omitempty" example:"summary"`
	Status    model.TaskStatus `json:"status,omitempty" example:"opened"`
}

type TaskResponse struct {
	Data TaskDto `json:"data"`
}

type TasksResponse struct {
	Pagination
	Data []TaskDto `json:"data"`
}

type CreateTaskDto struct {
	Summary string `json:"summary" binding:"required,min=1,max=2500" example:"summary"`
}
