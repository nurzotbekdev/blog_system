package controllers

import (
	"blog_system/config"
	"blog_system/helper"
	"blog_system/logging"
	"blog_system/models"
	"blog_system/schemas"
	"blog_system/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubscriptionController struct {
	SubscriptionService services.SubscriptionService
}

func NewSubscriptionController(subscription services.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{SubscriptionService: subscription}
}

func (subscription *SubscriptionController) Create(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create subscription attempt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	var body schemas.SubscriptionSchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid JSON format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	newSubscription := models.Subscription{
		UserID:    currentUserID,
		ChannelID: body.ChannelID,
	}

	if err := subscription.SubscriptionService.CreateSubscription(newSubscription); err != nil {
		if errors.Is(err, services.ErrChannelNotFound) {
			logging.Log.Warn("Channel not found for update", zap.Uint("channel_id", body.ChannelID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrAlreadySubscribed) {
			logging.Log.Warn("Duplicate subscription", zap.Uint("user_id", currentUserID), zap.Uint("channel_id", body.ChannelID))
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Update subscription failed", zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Subscription created successfully")
	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription created"})
}

func (subscription *SubscriptionController) SubscribedChannels(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create subscription attempt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	subscriptionData, err := subscription.SubscriptionService.GetSubscribedChannels(currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrSubscribedNotFound) {
			logging.Log.Warn("Subscribed not found", zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Get my subscribed failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": subscriptionData,
	})
}

func (subscription *SubscriptionController) ChannelSubscribers(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create subscription attempt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	var channel models.Channel
	err := config.DB.Where("user_id = ?", currentUserID).First(&channel).Error
	if err != nil {
		logging.Log.Warn("Channel not found", zap.Uint("user_id", currentUserID))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	subscriptionData, err := subscription.SubscriptionService.GetChannelSubscribers(channel.ID)
	if err != nil {
		logging.Log.Error("Get my subscribed failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": subscriptionData,
	})
}

func (subscription *SubscriptionController) RemoveSubscribers(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create subscription attempt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	subscriptionID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid subscription id param", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := subscription.SubscriptionService.DeleteSubscription(currentUserID, subscriptionID); err != nil {
		if errors.Is(err, services.ErrSubscribedNotFound) {
			logging.Log.Warn("Subscription not found", zap.Uint("user_id", currentUserID), zap.Uint("subscription_id", subscriptionID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}

		logging.Log.Error("Failed to delete subscription", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Uint("subscription_id", subscriptionID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Subscription deleted successfully", zap.Uint("user_id", currentUserID), zap.Uint("subscription_id", subscriptionID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription deleted"})
}

func (subscription *SubscriptionController) SubscriberStatistic(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create subscription attempt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	var channel models.Channel
	err := config.DB.Where("user_id = ?", currentUserID).First(&channel).Error
	if err != nil {
		logging.Log.Warn("Channel not found", zap.Uint("user_id", currentUserID))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	subscriptionStats, err := subscription.SubscriptionService.GetSubscriberStatistics(channel.ID)
	if err != nil {
		logging.Log.Error("Get my subscribed failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Subscriber statistics fetched successfully", zap.Uint("user_id", currentUserID), zap.Uint("channel_id", channel.ID))
	ctx.JSON(http.StatusOK, subscriptionStats)
}

func (subscription *SubscriptionController) SubscriberStatus(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create subscription attempt")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	channelID, err := helper.ParseUintParam(ctx, "channel_id")
	if err != nil {
		logging.Log.Warn("Invalid subscription id param", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscribed, err := subscription.SubscriptionService.IsSubscribed(currentUserID, channelID)
	if err != nil {
		logging.Log.Error("Get my subscribed failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"subscribed": subscribed})
}
