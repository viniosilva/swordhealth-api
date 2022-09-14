package model

import "time"

type TaskStatus string

const (
	TaskStatusOpened TaskStatus = "opened"
	TaskStatusClosed TaskStatus = "closed"
)

type Task struct {
	ID        int        `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	UserID int `db:"user_id"`

	Summary string     `db:"summary"`
	Status  TaskStatus `db:"status"`
}
