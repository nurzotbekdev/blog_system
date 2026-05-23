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

type PlaylistController struct {
	PlaylistService services.PlaylistService
}

func NewPlaylistController(playlist services.PlaylistService) *PlaylistController {
	return &PlaylistController{PlaylistService: playlist}
}

func (playlist *PlaylistController) Create(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	var body schemas.PlaylistRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Playlist created successfully"})
}

func (playlist *PlaylistController) GetMyPlaylist(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	playlistData, err := playlist.PlaylistService.GetMyPlaylist(currentUserID)
	if err != nil {
		if errors.Is(err, services.ErrPlaylistNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": playlistData})
}

func (playlist *PlaylistController) GetPlaylist(ctx *gin.Context) {
	playlistData, err := playlist.PlaylistService.GetPlaylist()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": playlistData})
}

func (playlist *PlaylistController) UpdatePlaylist(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	playlistID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.PlaylistUpdate
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := playlist.PlaylistService.EditePlaylist(currentUserID, playlistID, body); err != nil {
		if errors.Is(err, services.ErrPlaylistNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrNoFieldsToUpdate) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": services.ErrNoFieldsToUpdate.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Playlist update successful"})
}

func (playlist *PlaylistController) DeletePlaylist(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	playlistID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := playlist.PlaylistService.DeletePlaylist(currentUserID, playlistID); err != nil {
		if errors.Is(err, services.ErrPlaylistNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Playlist delete successful"})
}
