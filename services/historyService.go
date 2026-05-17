package services

import (
	"blog_system/config"
	"blog_system/models"
	"blog_system/schemas"
	"errors"
)

type HistoryService interface {
	GetUserHistory(userID uint, page, limit int) ([]schemas.HistoryResponse, error)
	DeleteHistory(userID, ID uint) error
}

type historyService struct{}

func NewHistoryServices() HistoryService {
	return &historyService{}
}

var (
	ErrHistoryNotFound = errors.New("History not found")
)

func (s *historyService) GetUserHistory(userID uint, page, limit int) ([]schemas.HistoryResponse, error) {
	var results []schemas.HistoryResponse

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}

	offset := (page - 1) * limit

	tx := config.DB.Table("histories").Select(`
		histories.id,
		histories.video_id,
		videos.title,
		videos.thumbnail_path,
		videos.duration_video,
		channels.name as channel_name,
		channels.profile_image,
		histories.created_at
	`).
		Joins("JOIN videos ON videos.id = histories.video_id").
		Joins("JOIN channels ON channels.id = videos.channel_id").
		Where("histories.user_id = ?", userID).
		Order("histories.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&results)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return results, nil
}

func (s *historyService) DeleteHistory(userID, ID uint) error {
	var history models.History
	if err := config.DB.Where("user_id = ? AND id = ?", userID, ID).First(&history).Error; err != nil {
		return ErrHistoryNotFound
	}

	return config.DB.Unscoped().Delete(&history).Error
}
