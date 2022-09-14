package service

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/viniosilva/swordhealth-api/internal/model"
	"github.com/viniosilva/swordhealth-api/internal/repository"
)

//go:generate mockgen -destination=../../mock/notification_service_mock.go -package=mock . NotificationService
type NotificationService interface {
	NotifyAdminUserOnSaveTask(ctx context.Context, task *model.Task) error
}

type notificationService struct {
	userRepository repository.UserRepository
	summaryLength  int
}

func NewNotificationService(userRepository repository.UserRepository, summaryLength int) NotificationService {
	return &notificationService{
		userRepository: userRepository,
		summaryLength:  summaryLength,
	}
}

func (impl *notificationService) NotifyAdminUserOnSaveTask(ctx context.Context, task *model.Task) error {
	users, _, err := impl.userRepository.ListUsers(ctx, 0, 0, repository.SetWhere("WHERE role = ?", []interface{}{model.UserRoleManager}))
	if err != nil {
		log.WithFields(log.Fields{
			"trace": "internal.service.notification.notifyadminuseronsavetask",
		}).Error(err.Error())
		return err
	}

	summary := task.Summary
	if len(summary) > impl.summaryLength {
		summary = summary[:impl.summaryLength] + "..."
	}

	for _, u := range users {
		log.WithFields(log.Fields{
			"trace": "internal.service.notification.notifyadminuseronsavetask",
			"user": map[string]interface{}{
				"id":       u.ID,
				"username": u.Username,
			},
			"task": map[string]interface{}{
				"id":         task.ID,
				"created_at": task.CreatedAt.Format("2006-01-02 15:04:05"),
				"updated_at": task.UpdatedAt.Format("2006-01-02 15:04:05"),
				"summary":    summary,
				"status":     task.Status,
			},
		}).Info("notification to manager user")
	}

	return nil
}
