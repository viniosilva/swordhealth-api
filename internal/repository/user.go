package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/model"
)

//go:generate mockgen -destination=../../mock/user_repository_mock.go -package=mock . UserRepository
type UserRepository interface {
	CreateUser(ctx context.Context, data dto.CreateUserDto) (*model.User, error)
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
