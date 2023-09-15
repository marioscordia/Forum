package service

import "newforum/internal/store"

type NotificationService struct{
	Notification store.NotificationStore
}

func NewNotificationService(store *store.Store) *NotificationService {
	return &NotificationService{
		Notification: store.NotificationStore,
	}
}

func (s *NotificationService) Notifications(userID int) ([]*store.Notification, error) {
	return s.Notification.Notifications(userID)
}

func (s *NotificationService) Update(userID int) error {
	return s.Notification.Update(userID)
}

func (s *NotificationService) NotificationNum(userID int) (int, error) {
	return s.Notification.NotificationNum(userID)
}
