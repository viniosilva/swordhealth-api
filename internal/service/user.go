package service

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/exception"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/repository"
)

//go:generate mockgen -destination=../../mock/user_service_mock.go -package=mock . UserService
type UserService interface {
	CreateUser(ctx context.Context, data dto.CreateUserDto) (*model.User, error)
	GetUserByID(ctx context.Context, id int) (*model.User, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (*model.User, error)
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
		log.WithContext(ctx).WithFields(log.Fields{
			"trace": "internal.service.user.createuser",
		}).Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (impl *userService) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	user, err := impl.userRepository.GetUserByID(ctx, id)
	if err != nil {
		if _, ok := err.(*exception.NotFoundException); !ok {
			log.WithContext(ctx).WithFields(log.Fields{
				"trace": "internal.service.user.getuserbyid",
			}).Error(err.Error())
		}

		return nil, err
	}

	return user, nil
}

func (impl *userService) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (*model.User, error) {
	user, err := impl.userRepository.GetUserByUsernameAndPassword(ctx, username, password)
	if err != nil {
		if _, ok := err.(*exception.NotFoundException); !ok {
			log.WithContext(ctx).WithFields(log.Fields{
				"trace": "internal.service.user.getuserbyusernameandpassword",
			}).Error(err.Error())
		}

		return nil, err
	}

	return user, nil
}
