package repository

import (
	"bytes"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/model"
)

//go:generate mockgen -destination=../../mock/user_repository_mock.go -package=mock . UserRepository
type UserRepository interface {
	CreateUser(ctx context.Context, data dto.CreateUserDto) (*model.User, error)
	ListUsers(ctx context.Context, limit, offset int, opts ...WhereOpt) ([]model.User, int, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (impl *userRepository) CreateUser(ctx context.Context, data dto.CreateUserDto) (*model.User, error) {
	now := time.Now()

	res, err := impl.db.ExecContext(ctx, `INSERT INTO users
			(created_at, updated_at, username, email, password, role)
			VALUES (?, ?, ?, ?, ?, ?);`,
		now, now, data.Username, data.Email, data.Password, data.Role)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        int(id),
		CreatedAt: now,
		UpdatedAt: now,
		Username:  data.Username,
		Email:     data.Email,
		Role:      data.Role,
	}, nil
}

func (impl *userRepository) ListUsers(ctx context.Context, limit, offset int, opts ...WhereOpt) ([]model.User, int, error) {
	var users []model.User
	total := 0

	var query bytes.Buffer
	query.WriteString(`
		SELECT id,
			created_at,
			updated_at,
			username,
			email,
			password,
			role
		FROM users
	`)

	args := []interface{}{}
	if len(opts) > 0 {
		query.WriteString(opts[0].Query())
		args = append(args, opts[0].Values()...)
	}
	if limit > 0 {
		query.WriteString("\nLIMIT ?")
		args = append(args, limit)
	}
	if offset > 0 {
		query.WriteString("\nOFFSET ?")
		args = append(args, offset)
	}

	err := impl.db.SelectContext(ctx, &users, query.String(), args...)
	if err != nil {
		return users, total, err
	}

	row := impl.db.QueryRowContext(ctx, `
		SELECT COUNT(id) as total
		FROM users
	`)
	err = row.Err()
	row.Scan(&total)

	return users, total, err
}
