package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/service"
	"github.com/viniosilva/swordhealth-api/mock"
)

func TestNotificationServiceNotifyAdminUserOnSaveTask(t *testing.T) {
	now := time.Now()

	var cases = map[string]struct {
		injectSummaryLength int
		inputTask           *model.Task
		mocking             func(userRepository *mock.MockUserRepository)
		expectedErr         error
	}{
		"should notify manager users": {
			injectSummaryLength: 30,
			inputTask: &model.Task{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				Summary:   "summary",
				Status:    model.TaskStatusOpened,
			},
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().ListUsers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]model.User{
						{ID: 1, Username: "user 1"},
						{ID: 2, Username: "user 2"},
					}, 2, nil)
			},
		},
		"should notify manager users when summary is longer than 30 characters": {
			injectSummaryLength: 30,
			inputTask: &model.Task{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				Summary:   "summary summary summary summary summary summary",
				Status:    model.TaskStatusOpened,
			},
			mocking: func(userRepository *mock.MockUserRepository) {
				userRepository.EXPECT().ListUsers(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]model.User{
						{ID: 1, Username: "user 1"},
						{ID: 2, Username: "user 2"},
					}, 2, nil)
			},
		},
		"should throw error when list users": {
			injectSummaryLength: 30,
			inputTask: &model.Task{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				Summary:   "summary",
				Status:    model.TaskStatusOpened,
			},
			mocking: func(userRepository *mock.MockUserRepository) {
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
			notificationService := service.NewNotificationService(userRepositoryMock, cs.injectSummaryLength)

			cs.mocking(userRepositoryMock)

			// when
			err := notificationService.NotifyAdminUserOnSaveTask(ctx, cs.inputTask)

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
	notificationService := service.NewNotificationService(userRepositoryMock, 30)

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
		notificationService.NotifyAdminUserOnSaveTask(ctx, task)
	}
}
