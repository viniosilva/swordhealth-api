package service

import (
	"context"
	"fmt"

	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/repository"
)

//go:generate mockgen -destination=../../mock/user_service_mock.go -package=mock . UserService
type UserService interface {
	CreateUser(ctx context.Context, data dto.CreateUserDto) (*model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (impl *userService) CreateUser(ctx context.Context, data dto.CreateUserDto) (*model.User, error) {
	if data.Role == "" {
		data.Role = model.UserRoleTechnician
	}

	user, err := impl.userRepository.CreateUser(ctx, data)
	if err != nil {
		fmt.Println("internal.service.user.createuser.error: ", err.Error())
		return nil, err
	}

	return user, nil
}
