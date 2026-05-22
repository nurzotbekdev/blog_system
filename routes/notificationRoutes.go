package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func NotificationRoutes(r *gin.Engine) {
	notificationService := services.NewNotificationServices()
	notificationController := controllers.NewNotificationController(notificationService)

	r.GET("/notifications", middleware.AuthMiddleware(), notificationController.GetNotifications)
	r.PUT("/notifications/:id", middleware.AuthMiddleware(), notificationController.MarkAsRead)
	r.GET("/unread-count", middleware.AuthMiddleware(), notificationController.GetUnreadCount)
}
