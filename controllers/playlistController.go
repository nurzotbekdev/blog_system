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

type PlaylistController struct {
	PlaylistService services.PlaylistService
}

func NewPlaylistController(playlist services.PlaylistService) *PlaylistController {
	return &PlaylistController{PlaylistService: playlist}
}

func (playlist *PlaylistController) Create(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized playlist create attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	var body schemas.PlaylistRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid playlist create request body", zap.Uint("user_id", currentUserID), zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	newPlaylist := models.Playlist{
		UserID:      currentUserID,
		Title:       body.Title,
		Description: body.Description,
		IsPrivate:   body.IsPrivate,
	}

	if err := playlist.PlaylistService.CreatePlaylist(newPlaylist); err != nil {
		logging.Log.Error("Failed to create playlist", zap.Uint("user_id", currentUserID), zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Playlist created successfully", zap.Uint("user_id", currentUserID), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
	ctx.JSON(http.StatusCreated, gin.H{"message": "Playlist created successfully"})
}

func (playlist *PlaylistController) GetMyPlaylist(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized get playlist attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	playlistData, err := playlist.PlaylistService.GetMyPlaylist(currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrPlaylistNotFound) {
			logging.Log.Warn("Playlist not found", zap.Uint("user_id", currentUserID), zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to get user playlists", zap.Uint("user_id", currentUserID), zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("User playlists fetched successfully", zap.Uint("user_id", currentUserID), zap.Int("count", len(playlistData)), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
	ctx.JSON(http.StatusOK, gin.H{"data": playlistData})
}

func (playlist *PlaylistController) GetPlaylist(ctx *gin.Context) {
	playlistData, err := playlist.PlaylistService.GetPlaylist()
	if err != nil {
		logging.Log.Error("Failed to get playlists", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Playlists fetched successfully", zap.Int("count", len(playlistData)), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
	ctx.JSON(http.StatusOK, gin.H{"data": playlistData})
}

func (playlist *PlaylistController) UpdatePlaylist(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized playlist update attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	playlistID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid playlist id param", zap.Uint("user_id", currentUserID), zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.PlaylistUpdate
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid playlist update JSON format", zap.Uint("user_id", currentUserID), zap.Uint("playlist_id", playlistID), zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := playlist.PlaylistService.EditePlaylist(currentUserID, playlistID, body); err != nil {
		if errors.Is(err, services.ErrPlaylistNotFound) {
			logging.Log.Warn("Playlist not found for update", zap.Uint("user_id", currentUserID), zap.Uint("playlist_id", playlistID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrNoFieldsToUpdate) {
			logging.Log.Warn("No fields to update in playlist", zap.Uint("user_id", currentUserID), zap.Uint("playlist_id", playlistID))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": services.ErrNoFieldsToUpdate.Error()})
			return
		}

		logging.Log.Error("Failed to update playlist", zap.Uint("user_id", currentUserID), zap.Uint("playlist_id", playlistID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Playlist updated successfully", zap.Uint("user_id", currentUserID), zap.Uint("playlist_id", playlistID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Playlist update successful"})
}

func (playlist *PlaylistController) DeletePlaylist(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized playlist delete attempt", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	playlistID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid playlist id param", zap.Uint("user_id", currentUserID), zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := playlist.PlaylistService.DeletePlaylist(currentUserID, playlistID); err != nil {
		if errors.Is(err, services.ErrPlaylistNotFound) {
			logging.Log.Warn("Playlist not found for delete", zap.Uint("user_id", currentUserID), zap.Uint("playlist_id", playlistID), zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to delete playlist", zap.Uint("user_id", currentUserID), zap.Uint("playlist_id", playlistID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Playlist deleted successfully", zap.Uint("user_id", currentUserID), zap.Uint("playlist_id", playlistID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Playlist delete successful"})
}
