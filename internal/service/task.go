package service

import (
	"context"

	log "github.com/sirupsen/logrus"
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
		log.WithFields(log.Fields{"error": err.Error()}).Error("internal.service.task.createtask")
		return nil, err
	}

	return task, nil
}
