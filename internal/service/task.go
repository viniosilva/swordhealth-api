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
	ListTasks(ctx context.Context, limit, offset int) ([]model.Task, int, error)
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
		log.WithContext(ctx).WithFields(log.Fields{
			"trace": "internal.service.task.createtask",
		}).Error(err.Error())
	}

	return task, err
}

func (impl *taskService) ListTasks(ctx context.Context, limit, offset int) ([]model.Task, int, error) {
	tasks, total, err := impl.taskRepository.ListTasks(ctx, limit, offset)
	if err != nil {
		log.WithContext(ctx).WithFields(log.Fields{
			"trace": "internal.service.task.listtasks",
		}).Error(err.Error())
	}

	return tasks, total, err
}
