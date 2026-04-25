package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func SubscriptionRoutes(r *gin.Engine) {
	subscriptionService := services.NewSubscriptionService()
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	r.POST("/subscription", middleware.AuthMiddleware(), subscriptionController.Create)
	r.GET("/subscription/my", middleware.AuthMiddleware(), subscriptionController.ChannelSubscribers)
	r.GET("/subscription/channel", middleware.AuthMiddleware(), subscriptionController.SubscribedChannels)
	r.DELETE("/subscription/:id", middleware.AuthMiddleware(), subscriptionController.RemoveSubscribers)
	r.GET("/subscription/stats", middleware.AuthMiddleware(), subscriptionController.SubscriberStatistic)
	r.GET("/subscription/:channel_id", middleware.AuthMiddleware(), subscriptionController.SubscriberStatus)
}
