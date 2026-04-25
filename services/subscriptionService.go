package services

import (
	"blog_system/config"
	"blog_system/logging"
	"blog_system/models"
	"blog_system/schemas"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SubscriptionService interface {
	CreateSubscription(subscription models.Subscription) error
	GetSubscribedChannels(userID uint) ([]schemas.SubscriptionResponse, error)
	GetChannelSubscribers(channelID uint) ([]schemas.ChannelSubscriptionResponse, error)
	DeleteSubscription(userID, ID uint) error
	GetSubscriberStatistics(channelID uint) (schemas.SubscriptionStatus, error)
	IsSubscribed(userID, channelID uint) (bool, error)
}

type subscriptionService struct{}

func NewSubscriptionService() SubscriptionService {
	return &subscriptionService{}
}

var (
	ErrAlreadySubscribed  = errors.New("user already subscribed to this channel")
	ErrSubscribedNotFound = errors.New("Subscribed not found")
)

func (s *subscriptionService) CreateSubscription(subscription models.Subscription) error {
	var channel models.Channel
	if err := config.DB.
		First(&channel, subscription.ChannelID).Error; err != nil {
		logging.Log.Warn("Channel not found", zap.Uint("channel_id", subscription.ChannelID), zap.Error(err))
		return ErrChannelNotFound
	}

	var existing models.Subscription
	err := config.DB.
		Where("user_id = ? AND channel_id = ?", subscription.UserID, subscription.ChannelID).
		First(&existing).Error
	if err == nil {
		logging.Log.Warn("Subscription already exists", zap.Uint("user_id", subscription.UserID), zap.Uint("channel_id", subscription.ChannelID))
		return ErrAlreadySubscribed
	}

	if err := config.DB.Create(&subscription).Error; err != nil {
		logging.Log.Error("Failed to create subscription", zap.Uint("user_id", subscription.UserID), zap.Uint("channel_id", subscription.ChannelID), zap.Error(err))
		return err
	}

	logging.Log.Info("Subscription created successfully", zap.Uint("user_id", subscription.UserID), zap.Uint("channel_id", subscription.ChannelID))
	if err := config.DB.
		Model(&models.Channel{}).
		Where("id = ?", subscription.ChannelID).
		UpdateColumn("total_subscribers",
			gorm.Expr("total_subscribers + ?", 1)).Error; err != nil {
		logging.Log.Error("Failed to update subscriber count", zap.Uint("channel_id", subscription.ChannelID), zap.Error(err))
		return err
	}

	return nil
}

func (s *subscriptionService) GetSubscribedChannels(userID uint) ([]schemas.SubscriptionResponse, error) {
	var results []schemas.SubscriptionResponse
	tx := config.DB.Table("subscriptions").Select(`
			subscriptions.id as id,
			channels.id as channel_id,
			channels.name,
			channels.description,
			channels.profile_image,
			channels.banner_image,
			channels.total_subscribers,
			channels.total_videos,
			channels.total_views`).
		Joins("JOIN channels ON channels.id = subscriptions.channel_id").
		Where("subscriptions.user_id = ?", userID).
		Scan(&results)

	if tx.Error != nil {
		logging.Log.Error("GetMySubscribed query failed", zap.Error(tx.Error), zap.Uint("user_id", userID))
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		logging.Log.Warn("Subscribed not found", zap.Uint("user_id", userID))
		return nil, ErrSubscribedNotFound
	}

	return results, nil
}

func (s *subscriptionService) GetChannelSubscribers(channelID uint) ([]schemas.ChannelSubscriptionResponse, error) {
	var results []schemas.ChannelSubscriptionResponse

	tx := config.DB.Table("subscriptions").
		Select(`
			subscriptions.id as id,
			users.id as user_id,
			users.full_name,
			users.profile_image
		`).
		Joins("JOIN users ON users.id = subscriptions.user_id").
		Where("subscriptions.channel_id = ?", channelID).
		Scan(&results)

	if tx.Error != nil {
		logging.Log.Error("GetChannelSubscribers query failed", zap.Error(tx.Error))
		return nil, tx.Error
	}

	return results, nil
}

func (s *subscriptionService) DeleteSubscription(userID, ID uint) error {
	var subscription models.Subscription
	if err := config.DB.
		Where("user_id = ? AND id = ?", userID, ID).
		First(&subscription).Error; err != nil {
		logging.Log.Warn("Subscription not found", zap.Uint("user_id", userID), zap.Uint("subscription_id", ID), zap.Error(err))

		return ErrSubscribedNotFound
	}

	if err := config.DB.Unscoped().Delete(&subscription).Error; err != nil {
		logging.Log.Error("Failed to delete subscription", zap.Uint("subscription_id", subscription.ID), zap.Error(err))
		return err
	}

	if err := config.DB.
		Model(&models.Channel{}).
		Where("id = ?", subscription.ChannelID).
		UpdateColumn("total_subscribers",
			gorm.Expr("GREATEST(total_subscribers-1,0)")).Error; err != nil {
		logging.Log.Error("Failed to update channel subscriber count", zap.Uint("channel_id", subscription.ChannelID), zap.Error(err))
		return err
	}

	logging.Log.Info("Channel subscriber count updated", zap.Uint("channel_id", subscription.ChannelID))
	return nil
}

func (s *subscriptionService) GetSubscriberStatistics(channelID uint) (schemas.SubscriptionStatus, error) {
	var status schemas.SubscriptionStatus
	now := time.Now()

	logging.Log.Info("Getting subscriber statistics started", zap.Uint("channel_id", channelID))

	if err := config.DB.Model(&models.Subscription{}).
		Where("channel_id = ?", channelID).
		Count(&status.Total).Error; err != nil {

		logging.Log.Error("Failed to get total subscribers", zap.Uint("channel_id", channelID), zap.Error(err))
		return status, err
	}

	if err := config.DB.Model(&models.Subscription{}).
		Where("channel_id = ? AND created_at >= ?", channelID, now.AddDate(0, 0, -1)).
		Count(&status.Today).Error; err != nil {

		logging.Log.Error("Failed to get today subscribers", zap.Uint("channel_id", channelID), zap.Error(err))
		return status, err
	}

	if err := config.DB.Model(&models.Subscription{}).
		Where("channel_id = ? AND created_at >= ?", channelID, now.AddDate(0, 0, -7)).
		Count(&status.ThisWeek).Error; err != nil {

		logging.Log.Error("Failed to get weekly subscribers", zap.Uint("channel_id", channelID), zap.Error(err))
		return status, err
	}

	if err := config.DB.Model(&models.Subscription{}).
		Where("channel_id = ? AND created_at >= ?", channelID, now.AddDate(0, -1, 0)).
		Count(&status.ThisMonth).Error; err != nil {

		logging.Log.Error("Failed to get monthly subscribers", zap.Uint("channel_id", channelID), zap.Error(err))
		return status, err
	}

	logging.Log.Info("Subscriber statistics fetched successfully",
		zap.Uint("channel_id", channelID),
		zap.Int64("total", status.Total),
		zap.Int64("today", status.Today),
		zap.Int64("week", status.ThisWeek),
		zap.Int64("month", status.ThisMonth),
	)

	return status, nil
}

func (s *subscriptionService) IsSubscribed(userID, channelID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Subscription{}).
		Where("user_id = ? AND channel_id = ?", userID, channelID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
