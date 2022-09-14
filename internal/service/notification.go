package service

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/repository"
)

//go:generate mockgen -destination=../../mock/notification_service_mock.go -package=mock . NotificationService
type NotificationService interface {
	NotifyAdminUserOnSaveTask(ctx context.Context, task *model.Task, actionUserID int) error
}

type notificationService struct {
	userRepository repository.UserRepository
}

func NewNotificationService(userRepository repository.UserRepository) NotificationService {
	return &notificationService{
		userRepository: userRepository,
	}
}

func (impl *notificationService) NotifyAdminUserOnSaveTask(ctx context.Context, task *model.Task, actionUserID int) error {
	actionUser, err := impl.userRepository.GetUserByID(ctx, actionUserID)
	if err != nil {
		log.WithFields(log.Fields{
			"trace": "internal.service.notification.notifyadminuseronsavetask",
		}).Error(err.Error())
		return err
	}

	if actionUser.Role == model.UserRoleManager {
		log.WithFields(log.Fields{
			"trace": "internal.service.notification.notifyadminuseronsavetask",
		}).Info("does not notify when user is manager")

		return nil
	}

	users, _, err := impl.userRepository.ListUsers(ctx, 0, 0, repository.SetWhere("WHERE role = ?", []interface{}{model.UserRoleManager}))
	if err != nil {
		log.WithFields(log.Fields{
			"trace": "internal.service.notification.notifyadminuseronsavetask",
		}).Error(err.Error())
		return err
	}

	message := fmt.Sprintf("the tech %s performed the task %d on date %s",
		actionUser.Username,
		task.ID,
		task.UpdatedAt.Format("2006-01-02 15:04:05"),
	)

	for _, u := range users {
		log.WithFields(log.Fields{
			"trace": "internal.service.notification.notifyadminuseronsavetask",
			"user": map[string]interface{}{
				"id":       u.ID,
				"username": u.Username,
			},
		}).Infof(message)
	}

	return nil
}
