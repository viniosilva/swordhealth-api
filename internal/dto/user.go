package dto

import (
	"github.com/viniosilva/swordhealth-api/internal/model"
)

type UserDto struct {
	ID        int    `json:"id" example:"1"`
	CreatedAt string `json:"created_at,omitempty" example:"1992-08-21 12:03:43"`
	UpdatedAt string `json:"updated_at,omitempty" example:"1992-08-21 12:03:43"`
	DeletedAt string `json:"deleted_at,omitempty" example:"1992-08-21 12:03:43"`

	Email    string         `json:"email,omitempty" example:"email@email.com"`
	Username string         `json:"username,omitempty" example:"username"`
	Role     model.UserRole `json:"role,omitempty" example:"technician"`
}

type UserResponse struct {
	Data UserDto `json:"data"`
}

type CreateUserDto struct {
	Username string         `json:"username" binding:"required,min=4,max=20" example:"username"`
	Email    string         `json:"email" binding:"required,email" example:"email@email.com"`
	Password string         `json:"password" binding:"required,min=4,max=20" example:"12345"`
	Role     model.UserRole `json:"role" binding:"enum" enums:"technician,manager" example:"technician"`
}
