package dto

import (
	"github.com/viniosilva/swordhealth-api/internal/model"
)

type UserDto struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`

	Email    string         `json:"email,omitempty"`
	Username string         `json:"username,omitempty"`
	Role     model.UserRole `json:"role,omitempty"`
}

type UserResponse struct {
	Data UserDto `json:"data"`
}

type CreateUserDto struct {
	Username string         `json:"username" binding:"required,min=4,max=20"`
	Email    string         `json:"email" binding:"required,email"`
	Password string         `json:"password" binding:"required,min=4,max=20"`
	Role     model.UserRole `json:"role" binding:"enum"`
}
