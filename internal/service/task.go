package service

import (
	"context"
	"fmt"

	"github.com/viniosilva/swordhealth-api/internal/dto"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/repository"
)

//go:generate mockgen -destination=../../mock/task_service_mock.go -package=mock . TaskService
type TaskService interface {
	CreateTask(ctx context.Context, data dto.CreateTaskDto) (*model.Task, error)
}

type taskService struct {
	taskRepository repository.TaskRepository
}

func NewTaskService(taskRepository repository.TaskRepository) TaskService {
	return &taskService{
		taskRepository: taskRepository,
	}
}

func (impl *taskService) CreateTask(ctx context.Context, data dto.CreateTaskDto) (*model.Task, error) {
	task, err := impl.taskRepository.CreateTask(ctx, data)
	if err != nil {
		fmt.Println("internal.service.task.createtask.error: ", err.Error())
		return nil, err
	}

	return task, nil
}
