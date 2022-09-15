package service

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/repository"
)

//go:generate mockgen -destination=../../mock/task_service_mock.go -package=mock . TaskService
type TaskService interface {
	CreateTask(ctx context.Context, userID int, summary string) (*model.Task, error)
	ListTasks(ctx context.Context, limit, offset int, user *model.User) ([]model.Task, int, error)
}

type taskService struct {
	taskRepository repository.TaskRepository
}

func NewTaskService(taskRepository repository.TaskRepository) TaskService {
	return &taskService{
		taskRepository: taskRepository,
	}
}

func (impl *taskService) CreateTask(ctx context.Context, userID int, summary string) (*model.Task, error) {
	task, err := impl.taskRepository.CreateTask(ctx, userID, summary)
	if err != nil {
		log.WithContext(ctx).WithFields(log.Fields{
			"trace": "internal.service.task.createtask",
		}).Error(err.Error())
	}

	return task, err
}

func (impl *taskService) ListTasks(ctx context.Context, limit, offset int, user *model.User) ([]model.Task, int, error) {
	opts := []repository.WhereOpt{}
	if user.Role != model.UserRoleManager {
		opts = append(opts, repository.SetWhere("WHERE user_id = ?", []interface{}{user.ID}))
	}

	tasks, total, err := impl.taskRepository.ListTasks(ctx, limit, offset, opts...)
	if err != nil {
		log.WithContext(ctx).WithFields(log.Fields{
			"trace": "internal.service.task.listtasks",
		}).Error(err.Error())
	}

	return tasks, total, err
}
