package controllers

import (
	"blog_system/helper"
	"blog_system/models"
	"blog_system/schemas"
	"blog_system/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	newPlaylistVideo := models.PlaylistVideo{
		PlaylistID: body.PlaylistID,
		VideoID:    body.VideoID,
	}

	if err := playlistVideo.PlaylistVideoService.CreatePlaylistVideo(newPlaylistVideo); err != nil {
		if errors.Is(err, services.ErrPlaylistNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrVideoNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrVideoAlreadyExistsInPlaylist) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Playlist video created successfully"})
}

func (playlistVideo *PlaylistVideoController) GetPlaylistDetail(ctx *gin.Context) {
	playlistVideoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playlistVideoData, err := playlistVideo.PlaylistVideoService.GetPlaylistDetail(playlistVideoID)
	if err != nil {
		if errors.Is(err, services.ErrPlaylistNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, playlistVideoData)
}

func (playlistVideo *PlaylistVideoController) GetMyPlaylistDetail(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	playlistVideoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	playlistVideoData, err := playlistVideo.PlaylistVideoService.GetMyPlaylistVideo(playlistVideoID, currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrPlaylistNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, playlistVideoData)
}

func (playlistVideo *PlaylistVideoController) DeletePlaylistVideo(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	playlistVideoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := playlistVideo.PlaylistVideoService.DeletePlaylistVideo(playlistVideoID, currentUserID); err != nil {
		if errors.Is(err, services.ErrPlaylistVideoNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrPlaylistNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Playlist video delete successful"})
}
