package services

import (
	"blog_system/config"
	"blog_system/models"
	"blog_system/schemas"
	"context"
	"errors"
	"fmt"
)

type NotificationService interface {
	GetNotifications(userID uint) ([]schemas.NotificationResponse, error)
	EditeMarkAsRead(notificationID, userID uint) error
	GetUnreadCount(userID uint) (int64, error)
}

type notificationService struct{}

func NewNotificationServices() NotificationService {
	return &notificationService{}
}

var (
	ErrNotificationsNotFound = errors.New("Notification not found")
)

func (s *notificationService) GetNotifications(userID uint) ([]schemas.NotificationResponse, error) {
	var notifications []models.Notification
	if err := config.DB.Preload("FromUser").Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, err
	}

	var response []schemas.NotificationResponse
	for _, n := range notifications {
		response = append(response, schemas.NotificationResponse{
			ID:         n.ID,
			UserID:     n.UserID,
			FromUserID: n.FromUserID,
			VideoID:    n.VideoID,
			CommentID:  n.CommentID,
			ChannelID:  n.ChannelID,
			CreatedAt:  n.CreatedAt,
		})
	}

	return response, nil
}

func (s *notificationService) EditeMarkAsRead(notificationID, userID uint) error {
	result := config.DB.
		Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)

	if result.RowsAffected == 0 {
		return ErrNotificationsNotFound
	}

	key := fmt.Sprintf("notif:unread:%d", userID)

	val, err := config.RedisClient.Get(context.Background(), key).Int()
	if err == nil && val > 0 {
		config.RedisClient.Decr(context.Background(), key)
	}

	return result.Error
}

func (s *notificationService) GetUnreadCount(userID uint) (int64, error) {
	key := fmt.Sprintf("notif:unread:%d", userID)

	val, err := config.RedisClient.Get(context.Background(), key).Int64()
	if err == nil {
		return val, nil
	}

	var count int64
	err = config.DB.
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error

	return count, err
}
