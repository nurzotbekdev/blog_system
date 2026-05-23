package controllers

import (
	"blog_system/config"
	"blog_system/helper"
	"blog_system/models"
	"blog_system/schemas"
	"blog_system/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	var body schemas.SubscriptionSchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	newSubscription := models.Subscription{
		UserID:    currentUserID,
		ChannelID: body.ChannelID,
	}

	if err := subscription.SubscriptionService.CreateSubscription(newSubscription); err != nil {
		if errors.Is(err, services.ErrChannelNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrAlreadySubscribed) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription created"})
}

func (subscription *SubscriptionController) SubscribedChannels(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	subscriptionData, err := subscription.SubscriptionService.GetSubscribedChannels(currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrSubscribedNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subscriptionData})
}

func (subscription *SubscriptionController) ChannelSubscribers(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	var channel models.Channel
	err := config.DB.Where("user_id = ?", currentUserID).First(&channel).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	subscriptionData, err := subscription.SubscriptionService.GetChannelSubscribers(channel.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": subscriptionData})
}

func (subscription *SubscriptionController) RemoveSubscribers(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	subscriptionID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := subscription.SubscriptionService.DeleteSubscription(currentUserID, subscriptionID); err != nil {
		if errors.Is(err, services.ErrSubscribedNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription deleted"})
}

func (subscription *SubscriptionController) SubscriberStatistic(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	var channel models.Channel
	err := config.DB.Where("user_id = ?", currentUserID).First(&channel).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	subscriptionStats, err := subscription.SubscriptionService.GetSubscriberStatistics(channel.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, subscriptionStats)
}

func (subscription *SubscriptionController) SubscriberStatus(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	channelID, err := helper.ParseUintParam(ctx, "channel_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscribed, err := subscription.SubscriptionService.IsSubscribed(currentUserID, channelID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"subscribed": subscribed})
}
