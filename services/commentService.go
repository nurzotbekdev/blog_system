package services

import (
	"blog_system/config"
	"blog_system/models"
	"blog_system/schemas"
	"errors"

	"gorm.io/gorm"
)

type CommentService interface {
	CreateComment(comment models.Comment) error
	GetVideoComments(page, limit int, videoID uint, sort string) (*schemas.CommentListResponse, error)
}

type commentService struct{}

func NewCommentServices() CommentService {
	return &commentService{}
}

var (
	ErrParentCommentNotFound = errors.New("parent comment not found")
)

func (s *commentService) CreateComment(comment models.Comment) error {
	var video models.Video
	if err := config.DB.First(&video, comment.VideoID).Error; err != nil {
		return ErrVideoNotFound
	}

	if comment.ParentID != nil {
		var parent models.Comment
		if err := config.DB.First(&parent, *comment.ParentID).Error; err != nil {
			return ErrParentCommentNotFound
		}
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		return err
	}

	if err := config.DB.Model(&models.Video{}).
		Where("id = ?", comment.VideoID).
		Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		return err
	}

	if err := config.DB.Model(&models.Channel{}).
		Where("id = ?", video.ChannelID).
		Update("total_comments", gorm.Expr("total_comments + 1")).Error; err != nil {
		return err
	}

	return nil
}

func (s *commentService) GetVideoComments(page, limit int, videoID uint, sort string) (*schemas.CommentListResponse, error) {
	var results []schemas.CommentResponse
	var total int64

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 50 {
		limit = 50
	}

	offset := (page - 1) * limit

	query := config.DB.Table("comments").
		Select(`
			comments.id,
			comments.user_id,
			users.full_name,
			users.profile_image as avatar,
			comments.content,
			comments.like_count,
			comments.created_at
		`).
		Joins("JOIN users ON users.id = comments.user_id").
		Where("comments.video_id = ? AND comments.parent_id IS NULL", videoID)

	switch sort {
	case "oldest":
		query = query.Order("comments.created_at ASC")
	case "popular":
		query = query.Order("comments.like_count DESC")
	default:
		query = query.Order("comments.created_at DESC")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	tx := query.
		Limit(limit).
		Offset(offset).
		Scan(&results)

	if tx.Error != nil {
		return nil, tx.Error
	}

	for i := range results {
		var replies []schemas.CommentResponse

		config.DB.Table("comments").
			Select(`
				comments.id,
				comments.user_id,
				users.full_name,
				users.profile_image as avatar,
				comments.content,
				comments.like_count,
				comments.created_at
			`).
			Joins("JOIN users ON users.id = comments.user_id").
			Where("comments.parent_id = ?", results[i].ID).
			Order("comments.created_at ASC").
			Scan(&replies)

		results[i].Replies = replies
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := schemas.CommentListResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Data:       results,
	}

	return &response, nil
}
