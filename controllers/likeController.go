package controllers

import (
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

type LikeController struct {
	LikeService services.LikeService
}

func NewLikeController(like services.LikeService) *LikeController {
	return &LikeController{LikeService: like}
}

func (like *LikeController) Create(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create video attempt")

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	currentUserID := userID.(uint)

	var body schemas.LikeRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid JSON format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	newLike := models.Like{
		UserID:  currentUserID,
		VideoID: body.VideoID,
		IsLike:  body.IsLike,
	}

	if err := like.LikeService.CreateLike(newLike); err != nil {
		if errors.Is(err, services.ErrVideoNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, services.ErrLikeAlready) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		logging.Log.Error("create video failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Like created successfully"})
}

func (like *LikeController) UpdateLike(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create video attempt")

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid subscription id param", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.LikeEdite
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid JSON format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := like.LikeService.EditeLike(currentUserID, videoID, body.IsLike); err != nil {
		if errors.Is(err, services.ErrLikeNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("create video failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Like update successfully"})
}

func (like *LikeController) DeleteLike(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create video attempt")

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid subscription id param", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := like.LikeService.DeleteLike(currentUserID, videoID); err != nil {
		if errors.Is(err, services.ErrLikeNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("create video failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Like delete successfully"})
}

func (like *LikeController) GetReaction(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create video attempt")

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid subscription id param", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	likeData, err := like.LikeService.GetReaction(currentUserID, videoID)
	if err != nil {
		if errors.Is(err, services.ErrLikeNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("create video failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, likeData)
}

func (like *LikeController) GetUserLikes(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create video attempt")

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	currentUserID := userID.(uint)

	likeData, err := like.LikeService.GetUserLikes(currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrLikeNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("create video failed")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": likeData})
}
