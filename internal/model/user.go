package model

import "time"

type UserRole string

const (
	UserRoleManager    UserRole = "manager"
	UserRoleTechnician UserRole = "technician"
)

type User struct {
	ID        int        `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	Email    string   `db:"email"`
	Username string   `db:"username"`
	Password string   `db:"password"`
	Role     UserRole `db:"role"`
}
