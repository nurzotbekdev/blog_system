package controllers

import (
	"blog_system/helper"
	"blog_system/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	NotificationService services.NotificationService
}

func NewNotificationController(notification services.NotificationService) *NotificationController {
	return &NotificationController{NotificationService: notification}
}

func (notification *NotificationController) GetNotifications(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	notificationData, err := notification.NotificationService.GetNotifications(currentUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": notificationData})
}

func (notification *NotificationController) MarkAsRead(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	notificationID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := notification.NotificationService.EditeMarkAsRead(notificationID, currentUserID); err != nil {
		if errors.Is(err, services.ErrNotificationsNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Marked as read"})
}

func (notification *NotificationController) GetUnreadCount(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	count, err := notification.NotificationService.GetUnreadCount(currentUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"unread": count})
}
