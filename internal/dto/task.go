package dto

import (
	"github.com/viniosilva/swordhealth-api/internal/model"
)

type TaskDto struct {
	ID        int              `json:"id"`
	CreatedAt string           `json:"created_at,omitempty"`
	UpdatedAt string           `json:"updated_at,omitempty"`
	User      UserDto          `json:"user,omitempty"`
	Summary   string           `json:"summary,omitempty"`
	Status    model.TaskStatus `json:"status,omitempty"`
}

type TaskResponse struct {
	Data TaskDto `json:"data"`
}

type CreateTaskDto struct {
	UserID  int    `json:"user_id" binding:"required,min=1"`
	Summary string `json:"summary" binding:"required,min=1,max=2500"`
}
