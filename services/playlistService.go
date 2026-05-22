package services

import (
	"blog_system/config"
	"blog_system/models"
	"blog_system/schemas"
	"errors"
)

type PlaylistService interface {
	CreatePlaylist(playlist models.Playlist) error
	GetMyPlaylist(userID uint) ([]schemas.PlaylistResponse, error)
	GetPlaylist() ([]schemas.PlaylistResponse, error)
	EditePlaylist(userID, ID uint, data schemas.PlaylistUpdate) error
	DeletePlaylist(userID, ID uint) error
}

type playlistService struct{}

func NewPlaylistServices() PlaylistService {
	return &playlistService{}
}

var (
	ErrPlaylistNotFound = errors.New("Playlist not found")
)

func (s *playlistService) CreatePlaylist(playlist models.Playlist) error {
	return config.DB.Create(&playlist).Error
}

func (s *playlistService) GetMyPlaylist(userID uint) ([]schemas.PlaylistResponse, error) {
	var playlist []models.Playlist
	if err := config.DB.Where("user_id = ?", userID).Find(&playlist).Error; err != nil {
		return nil, err
	}

	if len(playlist) == 0 {
		return nil, ErrPlaylistNotFound
	}

	var response []schemas.PlaylistResponse
	for _, p := range playlist {
		response = append(response, schemas.PlaylistResponse{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			IsPrivate:   p.IsPrivate,
			Thumbnail:   p.Thumbnail,
			VideoCount:  p.VideoCount,
			CreatedAt:   p.CreatedAt,
		})
	}

	return response, nil
}

func (s *playlistService) GetPlaylist() ([]schemas.PlaylistResponse, error) {
	var playlist []models.Playlist
	if err := config.DB.Where("is_private = ?", true).Find(&playlist).Error; err != nil {
		return nil, err
	}

	var response []schemas.PlaylistResponse
	for _, p := range playlist {
		response = append(response, schemas.PlaylistResponse{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			IsPrivate:   p.IsPrivate,
			Thumbnail:   p.Thumbnail,
			VideoCount:  p.VideoCount,
			CreatedAt:   p.CreatedAt,
		})
	}

	return response, nil
}

func (s *playlistService) EditePlaylist(userID, ID uint, data schemas.PlaylistUpdate) error {
	var playlist models.Playlist

	if err := config.DB.
		Where("user_id = ? AND id = ?", userID, ID).
		First(&playlist).Error; err != nil {
		return ErrPlaylistNotFound
	}

	updates := map[string]interface{}{}

	if data.Title != nil {
		updates["title"] = *data.Title
	}
	if data.Description != nil {
		updates["description"] = *data.Description
	}
	if data.IsPrivate != nil {
		updates["is_private"] = *data.IsPrivate
	}

	if len(updates) == 0 {
		return ErrNoFieldsToUpdate
	}

	if err := config.DB.Model(&playlist).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (s *playlistService) DeletePlaylist(userID, ID uint) error {
	var playlist models.Playlist

	if err := config.DB.
		Where("user_id = ? AND id = ?", userID, ID).
		First(&playlist).Error; err != nil {
		return ErrPlaylistNotFound
	}

	if err := config.DB.Delete(&playlist).Error; err != nil {
		return err
	}

	return nil
}
