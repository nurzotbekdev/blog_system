package controllers

import (
	"blog_system/helper"
	"blog_system/logging"
	"blog_system/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		logging.Log.Warn("Unauthorized get notifications attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	notificationData, err := notification.NotificationService.GetNotifications(currentUserID)
	if err != nil {
		logging.Log.Error("Failed to fetch notifications", zap.Error(err), zap.Uint("user_id", currentUserID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Notifications fetched successfully", zap.Uint("user_id", currentUserID), zap.Int("count", len(notificationData)))
	ctx.JSON(http.StatusOK, gin.H{"data": notificationData})
}

func (notification *NotificationController) MarkAsRead(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized mark notification as read attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	notificationID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid notification id parameter", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := notification.NotificationService.EditeMarkAsRead(notificationID, currentUserID); err != nil {
		if errors.Is(err, services.ErrNotificationsNotFound) {
			logging.Log.Warn("Notification not found for mark as read", zap.Uint("user_id", currentUserID), zap.Uint("notification_id", notificationID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}

		logging.Log.Error("Failed to mark notification as read", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Uint("notification_id", notificationID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Notification marked as read successfully", zap.Uint("user_id", currentUserID), zap.Uint("notification_id", notificationID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Marked as read"})
}

func (notification *NotificationController) GetUnreadCount(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized get unread notifications count attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	count, err := notification.NotificationService.GetUnreadCount(currentUserID)
	if err != nil {
		logging.Log.Error("Failed to fetch unread notifications count", zap.Error(err), zap.Uint("user_id", currentUserID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Unread notifications count fetched successfully", zap.Uint("user_id", currentUserID), zap.Int64("unread_count", count))
	ctx.JSON(http.StatusOK, gin.H{"unread": count})
}
