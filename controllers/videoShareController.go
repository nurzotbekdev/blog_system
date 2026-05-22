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

type VideoShareController struct {
	VideoShareService services.VideoShareService
}

func NewVideoShareController(share services.VideoShareService) *VideoShareController {
	return &VideoShareController{VideoShareService: share}
}

func (share *VideoShareController) Create(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized video share attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid video id param for share", zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.VideoShareSchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid video share request body", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	logging.Log.Info("Video share request started", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.String("platform", body.Platform))
	newVideoShare := models.VideoShare{
		UserID:   currentUserID,
		VideoID:  videoID,
		Platform: body.Platform,
	}

	if err := share.VideoShareService.CreateVideoShare(newVideoShare); err != nil {
		if errors.Is(err, services.ErrVideoNotFound) {
			logging.Log.Warn("Video not found for share", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrVideoNotVisibility) {
			logging.Log.Warn("Video is not visible for sharing", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to create video share", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Video shared successfully", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
	ctx.JSON(http.StatusCreated, gin.H{"message": "Video share successful"})
}
