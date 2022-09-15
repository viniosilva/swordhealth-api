package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/exception"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/service"
	"github.com/viniosilva/swordhealth-api/mock"
)

func TestUserServiceCreateUser(t *testing.T) {
	now := time.Now()
	user := &model.User{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		Username:  "username",
		Email:     "email@email.com",
		Password:  "aabbccddee",
		Role:      model.UserRoleTechnician,
	}

	var cases = map[string]struct {
		inputData    dto.CreateUserDto
		mocking      func(userRepository *mock.MockUserRepository)
		expectedUser *model.User
		expectedErr  error
	}{
		"should create user": {
			inputData: dto.CreateUserDto{
				Username: user.Username,
				Email:    user.Email,
				Password: "1122334455",
				Role:     user.Role,
			},
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(user, nil)
			},
			expectedUser: user,
		},
		"should create user when role is empty": {
			inputData: dto.CreateUserDto{
				Username: user.Username,
				Email:    user.Email,
				Password: "1122334455",
			},
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(user, nil)
			},
			expectedUser: user,
		},
		"should throw error when user repository create user": {
			inputData: dto.CreateUserDto{
				Username: user.Username,
				Email:    user.Email,
				Password: "1122334455",
				Role:     user.Role,
			},
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepositoryMock := mock.NewMockUserRepository(ctrl)
			userService := service.NewUserService(userRepositoryMock)

			cs.mocking(userRepositoryMock)

			// when
			user, err := userService.CreateUser(ctx, cs.inputData)

			// then
			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedUser, user)
		})
	}
}

func BenchmarkUserServiceCreateUser(b *testing.B) {
	// given
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	userService := service.NewUserService(userRepositoryMock)

	now := time.Now()
	user := &model.User{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		Username:  "username",
		Email:     "email@email.com",
		Password:  "aabbccddee",
		Role:      model.UserRoleManager,
	}

	userRepositoryMock.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
		AnyTimes().Return(user, nil)

	// when
	for i := 0; i < b.N; i++ {
		userService.CreateUser(ctx, dto.CreateUserDto{
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			Role:     user.Role,
		})
	}
}

func TestUserServiceGetUserByID(t *testing.T) {
	now := time.Now()
	user := &model.User{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		Username:  "username",
		Email:     "email@email.com",
		Password:  "aabbccddee",
		Role:      model.UserRoleTechnician,
	}

	var cases = map[string]struct {
		inputID      int
		mocking      func(userRepository *mock.MockUserRepository)
		expectedUser *model.User
		expectedErr  error
	}{
		"should return user": {
			inputID: user.ID,
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(user, nil)
			},
			expectedUser: user,
		},
		"should throw not found exception": {
			inputID: user.ID, mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(nil, &exception.NotFoundException{Message: "user not found"})
			},
			expectedErr: &exception.NotFoundException{Message: "user not found"},
		},
		"should throw error": {
			inputID: user.ID, mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepositoryMock := mock.NewMockUserRepository(ctrl)
			userService := service.NewUserService(userRepositoryMock)

			cs.mocking(userRepositoryMock)

			// when
			user, err := userService.GetUserByID(ctx, cs.inputID)

			// then
			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedUser, user)
		})
	}
}

func BenchmarkUserServiceGetUserByID(b *testing.B) {
	// given
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	userService := service.NewUserService(userRepositoryMock)

	now := time.Now()
	user := &model.User{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		Username:  "username",
		Email:     "email@email.com",
		Password:  "aabbccddee",
		Role:      model.UserRoleManager,
	}

	userRepositoryMock.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
		AnyTimes().Return(user, nil)

	// when
	for i := 0; i < b.N; i++ {
		userService.GetUserByID(ctx, user.ID)
	}
}

func TestUserServiceGetUserByUsernameAndPassword(t *testing.T) {
	now := time.Now()
	user := &model.User{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		Username:  "username",
		Email:     "email@email.com",
		Password:  "aabbccddee",
		Role:      model.UserRoleTechnician,
	}

	var cases = map[string]struct {
		inputUsername string
		inputPassword string
		mocking       func(userRepository *mock.MockUserRepository)
		expectedUser  *model.User
		expectedErr   error
	}{
		"should return user": {
			inputUsername: user.Username,
			inputPassword: user.Password,
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().GetUserByUsernameAndPassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil)
			},
			expectedUser: user,
		},
		"should throw not found exception": {
			inputUsername: user.Username,
			inputPassword: user.Password,
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().GetUserByUsernameAndPassword(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, &exception.NotFoundException{Message: "user not found"})
			},
			expectedErr: &exception.NotFoundException{Message: "user not found"},
		},
		"should throw error": {
			inputUsername: user.Username,
			inputPassword: user.Password,
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().GetUserByUsernameAndPassword(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepositoryMock := mock.NewMockUserRepository(ctrl)
			userService := service.NewUserService(userRepositoryMock)

			cs.mocking(userRepositoryMock)

			// when
			user, err := userService.GetUserByUsernameAndPassword(ctx, cs.inputUsername, cs.inputPassword)

			// then
			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedUser, user)
		})
	}
}

func BenchmarkUserServiceGetUserByUsernameAndPassword(b *testing.B) {
	// given
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	userService := service.NewUserService(userRepositoryMock)

	now := time.Now()
	user := &model.User{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		Username:  "username",
		Email:     "email@email.com",
		Password:  "aabbccddee",
		Role:      model.UserRoleManager,
	}

	userRepositoryMock.EXPECT().GetUserByUsernameAndPassword(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().Return(user, nil)

	// when
	for i := 0; i < b.N; i++ {
		userService.GetUserByUsernameAndPassword(ctx, user.Username, user.Password)
	}
}
