package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/swordhealth-api/internal/exception"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/service"
	"github.com/viniosilva/swordhealth-api/mock"
)

func TestNotificationServiceNotifyAdminUserOnSaveTask(t *testing.T) {
	now := time.Now()

	var cases = map[string]struct {
		inputTask         *model.Task
		inputActionUserID int
		mocking           func(userRepository *mock.MockUserRepository)
		expectedErr       error
	}{
		"should notify manager users": {
			inputTask: &model.Task{
				ID:        1,
				UserID:    1,
				CreatedAt: now,
				UpdatedAt: now,
				Summary:   "summary",
				Status:    model.TaskStatusOpened,
			},
			inputActionUserID: 1,
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(&model.User{ID: 1, Role: model.UserRoleTechnician}, nil)
				userRepository.EXPECT().ListUsers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]model.User{
						{ID: 2, Username: "user 2"},
						{ID: 3, Username: "user 3"},
					}, 2, nil)
			},
		},

		"should not notify when action user is a manager": {
			inputTask: &model.Task{
				ID:        1,
				UserID:    1,
				CreatedAt: now,
				UpdatedAt: now,
				Summary:   "summary",
				Status:    model.TaskStatusOpened,
			},
			inputActionUserID: 1,
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(&model.User{ID: 1, Role: model.UserRoleManager}, nil)
			},
		},
		"should throw not found exception when user not exist ": {
			inputTask: &model.Task{
				ID:        1,
				UserID:    1,
				CreatedAt: now,
				UpdatedAt: now,
				Summary:   "summary",
				Status:    model.TaskStatusOpened,
			},
			inputActionUserID: 1,
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(nil, &exception.NotFoundException{Message: "user not found"})
			},
			expectedErr: &exception.NotFoundException{Message: "user not found"},
		},
		"should throw error get user by id": {
			inputTask: &model.Task{
				ID:        1,
				UserID:    1,
				CreatedAt: now,
				UpdatedAt: now,
				Summary:   "summary",
				Status:    model.TaskStatusOpened,
			},
			inputActionUserID: 1,
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
		"should throw error when list users": {
			inputTask: &model.Task{
				ID:        1,
				UserID:    1,
				CreatedAt: now,
				UpdatedAt: now,
				Summary:   "summary",
				Status:    model.TaskStatusOpened,
			},
			inputActionUserID: 1,
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(&model.User{ID: 1, Role: model.UserRoleTechnician}, nil)
				userRepository.EXPECT().ListUsers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, 0, fmt.Errorf("error"))
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
			notificationService := service.NewNotificationService(userRepositoryMock)

			cs.mocking(userRepositoryMock)

			// when
			err := notificationService.NotifyAdminUserOnSaveTask(ctx, cs.inputTask, cs.inputActionUserID)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func BenchmarkNotificationServiceNotifyAdminUserOnSaveTask(b *testing.B) {
	// given
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	userRepositoryMock := mock.NewMockUserRepository(ctrl)
	notificationService := service.NewNotificationService(userRepositoryMock)

	userRepositoryMock.EXPECT().ListUsers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return([]model.User{
			{ID: 1, Username: "user 1"},
			{ID: 2, Username: "user 2"},
		}, 2, nil)

	now := time.Now()
	task := &model.Task{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    1,
		Summary:   "summary",
		Status:    model.TaskStatusOpened,
	}

	// when
	for i := 0; i < b.N; i++ {
		notificationService.NotifyAdminUserOnSaveTask(ctx, task, task.UserID)
	}
}
