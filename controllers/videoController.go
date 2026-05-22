package controllers

import (
	"blog_system/helper"
	"blog_system/logging"
	"blog_system/schemas"
	"blog_system/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VideoController struct {
	VideoService services.VideoService
}

func NewVideoController(video services.VideoService) *VideoController {
	return &VideoController{VideoService: video}
}

func (video *VideoController) Create(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create video attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	var req schemas.CreateVideoRequest
	if err := ctx.ShouldBind(&req); err != nil {
		logging.Log.Warn("Invalid video create request format", zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	filePath, _ := ctx.FormFile("file_path")
	thumbnailPath, _ := ctx.FormFile("thumbnail_path")

	req.FilePath = filePath
	req.ThumbnailPath = thumbnailPath

	if err := video.VideoService.CreateVideo(req, currentUserID); err != nil {
		if errors.Is(err, services.ErrChannelNotFound) {
			logging.Log.Warn("Channel not found for video creation", zap.Uint("user_id", currentUserID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrLanguageNotFound) {
			logging.Log.Warn("Language not found for video creation", zap.Uint("user_id", currentUserID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrCategoryNotFound) {
			logging.Log.Warn("Category not found for video creation", zap.Uint("user_id", currentUserID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to create video", zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Video created successfully", zap.Uint("user_id", currentUserID))
	ctx.JSON(http.StatusCreated, gin.H{"message": "Video created successfully"})
}

func (video *VideoController) MyVideo(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized access to my videos", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	videoData, err := video.VideoService.GetMyVideo(currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrChannelNotFound) {
			logging.Log.Warn("Channel not found for user videos", zap.Uint("user_id", currentUserID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrVideoNotFound) {
			logging.Log.Warn("Videos not found for user", zap.Uint("user_id", currentUserID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to fetch user videos", zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("User videos fetched successfully", zap.Uint("user_id", currentUserID), zap.Int("count", len(videoData)))
	ctx.JSON(http.StatusOK, gin.H{"data": videoData})
}

func (video *VideoController) ListVideos(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	categoryID := ctx.Query("category_id")
	search := ctx.Query("search")
	languageCode := ctx.Query("language_code")
	sortBy := ctx.Query("sort")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	category := 0
	if categoryID != "" {
		c, err := strconv.Atoi(categoryID)
		if err == nil {
			category = c
		}
	}

	videoData, err := video.VideoService.GetVideo(page, limit, uint(category), search, languageCode, sortBy)
	if err != nil {
		logging.Log.Error("Failed to list videos", zap.Error(err), zap.Int("page", page), zap.Int("limit", limit))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Videos fetched successfully", zap.Int("page", page), zap.Int("limit", limit))
	ctx.JSON(http.StatusOK, gin.H{"data": videoData})
}

func (video *VideoController) VideoDetail(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized access to video detail", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid video id param", zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	videoData, err := video.VideoService.GetVideoByID(videoID, currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrVideoNotFound) {
			logging.Log.Warn("Video not found", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to get video detail", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Video detail fetched successfully", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
	ctx.JSON(http.StatusOK, videoData)
}

func (video *VideoController) UpdateVideo(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized update video attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid video id param", zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	langIDStr := ctx.PostForm("language_id")
	catIDStr := ctx.PostForm("category_id")
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	visibility := ctx.PostForm("visibility")
	thumbnailPath, _ := ctx.FormFile("thumbnail_path")

	var languageID *uint
	var categoryID *uint

	if langIDStr != "" {
		id, err := strconv.ParseUint(langIDStr, 10, 64)
		if err != nil {
			logging.Log.Warn("Invalid language_id", zap.Uint("user_id", currentUserID), zap.String("language_id", langIDStr), zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid language_id"})
			return
		}
		val := uint(id)
		languageID = &val
	}

	if catIDStr != "" {
		id, err := strconv.ParseUint(catIDStr, 10, 64)
		if err != nil {
			logging.Log.Warn("Invalid category_id", zap.Uint("user_id", currentUserID), zap.String("category_id", catIDStr), zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
			return
		}
		val := uint(id)
		categoryID = &val
	}

	logging.Log.Info("Update video request started", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
	if err := video.VideoService.EditVideo(videoID, currentUserID, languageID, categoryID, &title, &description, &visibility, thumbnailPath); err != nil {
		if errors.Is(err, services.ErrChannelNotFound) {
			logging.Log.Warn("Channel not found for video update", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrVideoNotFound) {
			logging.Log.Warn("Video not found", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrLanguageNotFound) {
			logging.Log.Warn("Language not found", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrCategoryNotFound) {
			logging.Log.Warn("Category not found", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrNoFieldsToUpdate) {
			logging.Log.Warn("No fields to update", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to update video", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Video updated successfully", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Video updated successfully"})
}

func (video *VideoController) DeleteVideo(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized delete video attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid video id param", zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logging.Log.Info("delete video request started", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
	if err := video.VideoService.DeleteVideo(videoID, currentUserID); err != nil {
		if errors.Is(err, services.ErrChannelNotFound) {
			logging.Log.Warn("Channel not found for video delete", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrVideoNotFound) {
			logging.Log.Warn("Video not found for delete", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to delete video", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Video deleted successfully", zap.Uint("user_id", currentUserID), zap.Uint("video_id", videoID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Video deleted successfully"})
}
