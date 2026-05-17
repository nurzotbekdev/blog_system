package services

import (
	"blog_system/config"
	"blog_system/models"
	"blog_system/schemas"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LikeService interface {
	CreateLike(like models.Like) error
	EditeLike(userID, videoID uint, isLike bool) error
	DeleteLike(userID, videoID uint) error
	GetReaction(userID, videoID uint) (*schemas.LikeEdite, error)
	GetUserLikes(userID uint) ([]schemas.UserLikeResponse, error)
}

type likeService struct{}

func NewLikeServices() LikeService {
	return &likeService{}
}

var (
	ErrLikeAlready  = errors.New("reaction already exists")
	ErrLikeNotFound = errors.New("Like not found")
)

func (s *likeService) CreateLike(like models.Like) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var video models.Video
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&video, like.VideoID).Error; err != nil {
			return ErrVideoNotFound
		}

		var exists models.Like
		if err := tx.Where("user_id = ? AND video_id = ?", like.UserID, like.VideoID).First(&exists).Error; err == nil {
			return ErrLikeAlready
		}

		if err := tx.Create(&like).Error; err != nil {
			return err
		}

		if like.IsLike {
			if err := tx.Model(&models.Video{}).
				Where("id = ?", like.VideoID).
				UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&models.Video{}).
				Where("id = ?", like.VideoID).
				UpdateColumn("dislike_count", gorm.Expr("dislike_count + 1")).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *likeService) EditeLike(userID, videoID uint, isLike bool) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var like models.Like
		if err := tx.Where("user_id = ? AND video_id = ?", userID, videoID).First(&like).Error; err != nil {
			return ErrLikeNotFound
		}

		if like.IsLike {
			if err := tx.Model(&models.Video{}).
				Where("id = ?", like.VideoID).
				UpdateColumn("like_count",
					gorm.Expr("GREATEST(like_count - 1,0)")).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&models.Video{}).
				Where("id = ?", like.VideoID).
				UpdateColumn("dislike_count",
					gorm.Expr("GREATEST(dislike_count - 1,0)")).Error; err != nil {
				return err
			}
		}

		if isLike {
			if err := tx.Model(&models.Video{}).
				Where("id = ?", like.VideoID).
				UpdateColumn("like_count",
					gorm.Expr("like_count + 1")).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&models.Video{}).
				Where("id = ?", like.VideoID).
				UpdateColumn("dislike_count",
					gorm.Expr("dislike_count + 1")).Error; err != nil {
				return err
			}
		}

		like.IsLike = isLike
		return tx.Save(&like).Error
	})
}

func (s *likeService) DeleteLike(userID, videoID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var like models.Like
		if err := tx.
			Where("user_id = ? AND video_id = ?", userID, videoID).
			First(&like).Error; err != nil {
			return ErrLikeNotFound
		}

		if like.IsLike {
			if err := tx.Model(&models.Video{}).
				Where("id = ?", like.VideoID).
				UpdateColumn("like_count",
					gorm.Expr("GREATEST(like_count - 1,0)")).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&models.Video{}).
				Where("id = ?", like.VideoID).
				UpdateColumn("dislike_count",
					gorm.Expr("GREATEST(dislike_count - 1,0)")).Error; err != nil {
				return err
			}
		}

		if err := tx.Delete(&like).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *likeService) GetReaction(userID, videoID uint) (*schemas.LikeEdite, error) {
	var like models.Like
	if err := config.DB.Where("user_id = ? AND video_id = ?", userID, videoID).
		First(&like).Error; err != nil {
		return nil, ErrLikeNotFound
	}

	result := schemas.LikeEdite{
		IsLike: like.IsLike,
	}

	return &result, nil
}

func (s *likeService) GetUserLikes(userID uint) ([]schemas.UserLikeResponse, error) {
	var like models.Like
	if err := config.DB.Where("user_id = ?", userID).First(&like).Error; err != nil {
		return nil, ErrLikeNotFound
	}

	var results []schemas.UserLikeResponse
	tx := config.DB.Table("likes").Select(`
		likes.id,
		likes.video_id,
		videos.title,
		likes.is_like
	`).
		Joins("JOIN videos ON videos.id = like.video_id").
		Where("like.video_id = ?", like.ID).
		Scan(&results)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, ErrLikeNotFound
	}

	return results, nil
}
