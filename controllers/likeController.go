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
		logging.Log.Warn("Unauthorized request to create like", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	var body schemas.LikeRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid json format for like create", zap.Error(err), zap.Uint("user_id", currentUserID))
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
			logging.Log.Warn("Video not found while creating like", zap.Uint("user_id", currentUserID), zap.Uint("video_id", body.VideoID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrLikeAlready) {
			logging.Log.Warn("Duplicate like detected", zap.Uint("user_id", currentUserID), zap.Uint("video_id", body.VideoID), zap.Bool("is_like", body.IsLike))
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to create like", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Uint("video_id", body.VideoID), zap.Bool("is_like", body.IsLike))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Like created successfully", zap.Uint("user_id", currentUserID), zap.Uint("video_id", body.VideoID), zap.Bool("is_like", body.IsLike))
	ctx.JSON(http.StatusCreated, gin.H{"message": "Like created successfully"})
}

func (like *LikeController) UpdateLike(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized like update attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid video id parameter", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.LikeEdite
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid JSON format for like update", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := like.LikeService.EditeLike(currentUserID, videoID, body.IsLike); err != nil {
		if errors.Is(err, services.ErrLikeNotFound) {
			logging.Log.Warn("Like not found for update", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to update like", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Bool("is_like", body.IsLike))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Like updated successfully", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Bool("is_like", body.IsLike))
	ctx.JSON(http.StatusOK, gin.H{"message": "Like update successfully"})
}

func (like *LikeController) DeleteLike(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized like delete attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid video id parameter", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := like.LikeService.DeleteLike(currentUserID, videoID); err != nil {
		if errors.Is(err, services.ErrLikeNotFound) {
			logging.Log.Warn("Like not found for delete", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to delete like", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Like deleted successfully", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Like delete successfully"})
}

func (like *LikeController) GetReaction(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized get reaction attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid video id parameter", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	likeData, err := like.LikeService.GetReaction(currentUserID, videoID)
	if err != nil {
		if errors.Is(err, services.ErrLikeNotFound) {
			logging.Log.Warn("Like reaction not found", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to get like reaction", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Like reaction fetched successfully", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
	ctx.JSON(http.StatusOK, likeData)
}

func (like *LikeController) GetUserLikes(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized get user likes attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	likeData, err := like.LikeService.GetUserLikes(currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrLikeNotFound) {
			logging.Log.Warn("User likes not found", zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to fetch user likes", zap.Error(err), zap.Uint("user_id", currentUserID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("User likes fetched successfully", zap.Uint("user_id", currentUserID), zap.Int("count", len(likeData)))
	ctx.JSON(http.StatusOK, gin.H{"data": likeData})
}
