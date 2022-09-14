package repository

import (
	"bytes"
	"context"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/exception"
	"github.com/viniosilva/swordhealth-api/internal/model"
)

//go:generate mockgen -destination=../../mock/task_repository_mock.go -package=mock . TaskRepository
type TaskRepository interface {
	CreateTask(ctx context.Context, data dto.CreateTaskDto) (*model.Task, error)
	ListTasks(ctx context.Context, limit, offset int) ([]model.Task, int, error)
}

type taskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (impl *taskRepository) CreateTask(ctx context.Context, data dto.CreateTaskDto) (*model.Task, error) {
	now := time.Now()

	res, err := impl.db.ExecContext(ctx, `INSERT INTO tasks
			(created_at, updated_at, user_id, summary, status)
			VALUES (?, ?, ?, ?, ?);`,
		now, now, data.UserID, data.Summary, model.TaskStatusOpened)
	if err != nil {
		if e, ok := err.(*mysql.MySQLError); ok && int(e.Number) == int(MySQLErrorCodeForeignKeyConstraint) {
			err = &exception.ForeignKeyConstraintException{Message: "user not found"}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.Task{
		ID:        int(id),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    data.UserID,
		Summary:   data.Summary,
		Status:    model.TaskStatusOpened,
	}, nil
}

func (impl *taskRepository) ListTasks(ctx context.Context, limit, offset int) ([]model.Task, int, error) {
	var tasks []model.Task
	total := 0

	var query bytes.Buffer
	query.WriteString(`
		SELECT id,
			created_at,
			updated_at,
			deleted_at,
			user_id,
			summary,
			status
		FROM tasks
	`)

	args := []interface{}{}
	if limit > 0 {
		query.WriteString("\nLIMIT ?")
		args = append(args, limit)
	}
	if offset > 0 {
		query.WriteString("\nOFFSET ?")
		args = append(args, offset)
	}

	err := impl.db.SelectContext(ctx, &tasks, query.String(), args...)
	if err != nil {
		return tasks, total, err
	}

	row := impl.db.QueryRowContext(ctx, `
		SELECT COUNT(id) as total
		FROM tasks
	`)
	err = row.Err()
	row.Scan(&total)

	return tasks, total, err
}
