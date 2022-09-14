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

func TestTaskServiceCreateTask(t *testing.T) {
	now := time.Now()
	task := &model.Task{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    1,
		Summary:   "summary",
		Status:    model.TaskStatusOpened,
	}

	var cases = map[string]struct {
		inputData    dto.CreateTaskDto
		mocking      func(taskRepository *mock.MockTaskRepository)
		expectedTask *model.Task
		expectedErr  error
	}{
		"should create task": {
			inputData: dto.CreateTaskDto{
				UserID:  task.UserID,
				Summary: task.Summary,
			},
			mocking: func(taskRepository *mock.MockTaskRepository) {
				taskRepository.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(task, nil)
			},
			expectedTask: task,
		},
		"should throw foreign key constraint exception when user not found": {
			inputData: dto.CreateTaskDto{
				UserID:  task.UserID,
				Summary: task.Summary,
			},
			mocking: func(taskRepository *mock.MockTaskRepository) {
				taskRepository.EXPECT().CreateTask(gomock.Any(), gomock.Any()).
					Return(nil, &exception.ForeignKeyConstraintException{Message: "user not found"})
			},
			expectedErr: &exception.ForeignKeyConstraintException{Message: "user not found"},
		},
		"should throw error when task repository create task": {
			inputData: dto.CreateTaskDto{
				UserID:  task.UserID,
				Summary: task.Summary,
			},
			mocking: func(taskRepository *mock.MockTaskRepository) {
				taskRepository.EXPECT().CreateTask(gomock.Any(), gomock.Any()).
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

			taskRepositoryMock := mock.NewMockTaskRepository(ctrl)
			taskService := service.NewTaskService(taskRepositoryMock)

			cs.mocking(taskRepositoryMock)

			// when
			task, err := taskService.CreateTask(ctx, cs.inputData)

			// then
			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedTask, task)
		})
	}
}

func BenchmarkTaskServiceCreateTask(b *testing.B) {
	// given
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	taskRepositoryMock := mock.NewMockTaskRepository(ctrl)
	taskService := service.NewTaskService(taskRepositoryMock)

	now := time.Now()
	task := &model.Task{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    1,
		Summary:   "summary",
		Status:    model.TaskStatusOpened,
	}

	taskRepositoryMock.EXPECT().CreateTask(gomock.Any(), gomock.Any()).
		AnyTimes().Return(task, nil)

	// when
	for i := 0; i < b.N; i++ {
		taskService.CreateTask(ctx, dto.CreateTaskDto{
			UserID:  task.ID,
			Summary: task.Summary,
		})
	}
}

func TestTaskServiceListTasks(t *testing.T) {
	now := time.Now()
	task := model.Task{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    1,
		Summary:   "summary",
		Status:    model.TaskStatusOpened,
	}

	var cases = map[string]struct {
		inputLimit    int
		inputOffset   int
		mocking       func(taskRepository *mock.MockTaskRepository)
		expectedTasks []model.Task
		expectedTotal int
		expectedErr   error
	}{
		"should list tasks": {
			inputLimit:  10,
			inputOffset: 0,
			mocking: func(taskRepository *mock.MockTaskRepository) {
				taskRepository.EXPECT().ListTasks(gomock.Any(), gomock.Any(), gomock.Any()).Return([]model.Task{task}, 1, nil)
			},
			expectedTasks: []model.Task{task},
			expectedTotal: 1,
		},
		"should throw error when task repository create task": {
			inputLimit:  10,
			inputOffset: 0,
			mocking: func(taskRepository *mock.MockTaskRepository) {
				taskRepository.EXPECT().ListTasks(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, 0, fmt.Errorf("error"))
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

			taskRepositoryMock := mock.NewMockTaskRepository(ctrl)
			taskService := service.NewTaskService(taskRepositoryMock)

			cs.mocking(taskRepositoryMock)

			// when
			tasks, total, err := taskService.ListTasks(ctx, cs.inputLimit, cs.inputOffset)

			// then
			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedTasks, tasks)
			assert.Equal(t, cs.expectedTotal, total)
		})
	}
}

func BenchmarkTaskServiceListTasks(b *testing.B) {
	// given
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	taskRepositoryMock := mock.NewMockTaskRepository(ctrl)
	taskService := service.NewTaskService(taskRepositoryMock)

	now := time.Now()
	tasks := []model.Task{{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    1,
		Summary:   "summary",
		Status:    model.TaskStatusOpened,
	}}

	taskRepositoryMock.EXPECT().ListTasks(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().Return(tasks, 1, nil)

	// when
	for i := 0; i < b.N; i++ {
		taskService.ListTasks(ctx, 10, 0)
	}
}
