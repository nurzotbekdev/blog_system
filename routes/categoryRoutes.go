package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(r *gin.Engine) {
	categoryService := services.NewCategoryServices()
	categoryController := controllers.NewCategoryController(categoryService)

	r.POST("/category", middleware.AuthMiddleware(), middleware.AdminOnly(), categoryController.Create)
	r.GET("/category", middleware.AuthMiddleware(), categoryController.AllCategory)
	r.PUT("/category/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), categoryController.UpdateCategory)
	r.DELETE("/category/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), categoryController.RemoveCategory)
}
