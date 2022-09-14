package service

import (
	"context"

	log "github.com/sirupsen/logrus"
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
	cryptoService  CryptoService
}

func NewUserService(userRepository repository.UserRepository, cryptoService CryptoService) UserService {
	return &userService{
		userRepository: userRepository,
		cryptoService:  cryptoService,
	}
}

func (impl *userService) CreateUser(ctx context.Context, data dto.CreateUserDto) (*model.User, error) {
	if data.Role == "" {
		data.Role = model.UserRoleTechnician
	}

	data.Password = impl.cryptoService.Hash(data.Password)
	user, err := impl.userRepository.CreateUser(ctx, data)
	if err != nil {
		log.WithFields(log.Fields{
			"trace": "internal.service.user.createuser",
		}).Error(err.Error())
		return nil, err
	}

	return user, nil
}
