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

type PlaylistVideoController struct {
	PlaylistVideoService services.PlaylistVideoService
}

func NewPlaylistVideoController(playlistVideo services.PlaylistVideoService) *PlaylistVideoController {
	return &PlaylistVideoController{PlaylistVideoService: playlistVideo}
}

func (playlistVideo *PlaylistVideoController) Create(ctx *gin.Context) {
	var body schemas.PlaylistVideoRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid playlist video request body", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	newPlaylistVideo := models.PlaylistVideo{
		PlaylistID: body.PlaylistID,
		VideoID:    body.VideoID,
	}

	if err := playlistVideo.PlaylistVideoService.CreatePlaylistVideo(newPlaylistVideo); err != nil {
		if errors.Is(err, services.ErrPlaylistNotFound) {
			logging.Log.Warn("Playlist not found while adding video", zap.Uint("playlist_id", body.PlaylistID), zap.Uint("video_id", body.VideoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrVideoNotFound) {
			logging.Log.Warn("Video not found while adding to playlist", zap.Uint("playlist_id", body.PlaylistID), zap.Uint("video_id", body.VideoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrVideoAlreadyExistsInPlaylist) {
			logging.Log.Warn("Video already exists in playlist", zap.Uint("playlist_id", body.PlaylistID), zap.Uint("video_id", body.VideoID), zap.Error(err))
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to add video to playlist", zap.Uint("playlist_id", body.PlaylistID), zap.Uint("video_id", body.VideoID), zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Video added to playlist successfully", zap.Uint("playlist_id", body.PlaylistID), zap.Uint("video_id", body.VideoID))
	ctx.JSON(http.StatusCreated, gin.H{"message": "Playlist video created successfully"})
}

func (playlistVideo *PlaylistVideoController) GetPlaylistDetail(ctx *gin.Context) {
	playlistVideoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid playlist id param", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playlistVideoData, err := playlistVideo.PlaylistVideoService.GetPlaylistDetail(playlistVideoID)
	if err != nil {
		if errors.Is(err, services.ErrPlaylistNotFound) {
			logging.Log.Warn("Playlist not found", zap.Uint("playlist_id", playlistVideoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to get playlist detail", zap.Uint("playlist_id", playlistVideoID), zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Playlist detail fetched successfully", zap.Uint("playlist_id", playlistVideoID), zap.Int("video_count", len(playlistVideoData.Videos)))
	ctx.JSON(http.StatusOK, playlistVideoData)
}

func (playlistVideo *PlaylistVideoController) GetMyPlaylistDetail(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized playlist detail attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	playlistVideoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid playlist id param", zap.Uint("user_id", currentUserID), zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playlistVideoData, err := playlistVideo.PlaylistVideoService.GetMyPlaylistVideo(playlistVideoID, currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrPlaylistNotFound) {
			logging.Log.Warn("Playlist not found for user", zap.Uint("user_id", currentUserID), zap.Uint("playlist_id", playlistVideoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to get user playlist detail", zap.Uint("user_id", currentUserID), zap.Uint("playlist_id", playlistVideoID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("User playlist detail fetched successfully", zap.Uint("user_id", currentUserID), zap.Uint("playlist_id", playlistVideoID))
	ctx.JSON(http.StatusOK, playlistVideoData)
}

func (playlistVideo *PlaylistVideoController) DeletePlaylistVideo(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized playlist video delete attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	playlistVideoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid playlist video id param", zap.Uint("user_id", currentUserID), zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := playlistVideo.PlaylistVideoService.DeletePlaylistVideo(playlistVideoID, currentUserID); err != nil {
		if errors.Is(err, services.ErrPlaylistVideoNotFound) {
			logging.Log.Warn("Playlist video not found", zap.Uint("user_id", currentUserID), zap.Uint("playlist_video_id", playlistVideoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrPlaylistNotFound) {
			logging.Log.Warn("Playlist not found for playlist video delete", zap.Uint("user_id", currentUserID), zap.Uint("playlist_video_id", playlistVideoID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to delete playlist video", zap.Uint("user_id", currentUserID), zap.Uint("playlist_video_id", playlistVideoID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Playlist video deleted successfully", zap.Uint("user_id", currentUserID), zap.Uint("playlist_video_id", playlistVideoID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Playlist video delete successful"})
}
