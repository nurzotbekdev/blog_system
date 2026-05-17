package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func HistoryRoutes(r *gin.Engine) {
	historyService := services.NewHistoryServices()
	historyController := controllers.NewHistoryController(historyService)

	r.GET("/history/me", middleware.AuthMiddleware(), historyController.GetUserHistory)
	r.DELETE("/history/:id", middleware.AuthMiddleware(), historyController.DeleteHistory)
}
