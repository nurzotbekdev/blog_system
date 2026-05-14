package routes

import (
	"blog_system/controllers"
	"blog_system/middleware"
	"blog_system/services"

	"github.com/gin-gonic/gin"
)

func CommentReactionRoutes(r *gin.Engine) {
	commentReactionService := services.NewCommentRectionServices()
	commentReactionController := controllers.NewCommentRectionController(commentReactionService)

	r.POST("/comments/:id/reaction", middleware.AuthMiddleware(), commentReactionController.CreateCommentReaction)
}
