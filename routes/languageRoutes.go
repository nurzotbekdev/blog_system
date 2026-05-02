package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func LanguageRoutes(r *gin.Engine) {
	languageService := services.NewLanguageServices()
	languageController := controllers.NewLanguageController(languageService)

	r.POST("/language", middleware.AuthMiddleware(), middleware.AdminOnly(), languageController.Create)
	r.GET("/language", middleware.AuthMiddleware(), languageController.AllLanguage)
	r.PUT("/language/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), languageController.UpdateLanguage)
	r.DELETE("/language/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), languageController.RemoveLanguage)
}
