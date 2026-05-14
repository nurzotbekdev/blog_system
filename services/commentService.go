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
	EditVideoComment(userID, ID uint, content string) error
	DeleteComment(userID, ID uint) error
	CreateReply(parentCommentID uint, comment models.Comment) error
	GetStats(ID uint) (*schemas.CommentStatsResponse, error)
}

type commentService struct{}

func NewCommentServices() CommentService {
	return &commentService{}
}

var (
	ErrParentCommentNotFound = errors.New("parent comment not found")
	ErrCommentNotFound       = errors.New("Comment not found")
	ErrForbidden             = errors.New("You are not allowed to perform this action")
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

func (s *commentService) EditVideoComment(userID, ID uint, content string) error {
	var comment models.Comment
	if err := config.DB.Where("id = ? AND user_id = ?", ID, userID).First(&comment).Error; err != nil {
		return ErrCommentNotFound
	}

	if err := config.DB.Model(&comment).
		Update("content", content).Error; err != nil {
		return err
	}

	return nil
}

func (s *commentService) DeleteComment(userID, ID uint) error {
	var comment models.Comment

	if err := config.DB.First(&comment, ID).Error; err != nil {
		return ErrCommentNotFound
	}

	var video models.Video
	if err := config.DB.First(&video, comment.VideoID).Error; err != nil {
		return ErrVideoNotFound
	}

	var channel models.Channel
	if err := config.DB.First(&channel, video.ChannelID).Error; err != nil {
		return ErrChannelNotFound
	}

	isCommentOwner := comment.UserID == userID
	isVideoOwner := channel.UserID == userID

	if !isCommentOwner && !isVideoOwner {
		return ErrForbidden
	}

	if err := config.DB.Delete(&comment).Error; err != nil {
		return err
	}

	if err := config.DB.Model(&models.Video{}).Where("id = ?", comment.VideoID).UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1,0)")).Error; err != nil {
		return err
	}

	if err := config.DB.Model(&models.Channel{}).Where("id = ?", video.ChannelID).UpdateColumn("total_comments", gorm.Expr("GREATEST(total_comments - 1,0)")).Error; err != nil {
		return err
	}
	return nil
}

func (s *commentService) CreateReply(parentCommentID uint, comment models.Comment) error {
	var parent models.Comment
	if err := config.DB.First(&parent, parentCommentID).Error; err != nil {
		return ErrParentCommentNotFound
	}

	comment.VideoID = parent.VideoID
	comment.ParentID = &parentCommentID

	if err := config.DB.Create(&comment).Error; err != nil {
		return err
	}

	if err := config.DB.Model(&models.Video{}).
		Where("id = ?", parent.VideoID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		return err
	}

	if err := config.DB.Model(&models.Channel{}).
		Where("id = (?)",
			config.DB.Table("videos").
				Select("channel_id").
				Where("id = ?", parent.VideoID),
		).
		UpdateColumn("total_comments", gorm.Expr("total_comments + 1")).Error; err != nil {
		return err
	}

	return nil
}

func (s *commentService) GetStats(ID uint) (*schemas.CommentStatsResponse, error) {
	var comment models.Comment
	if err := config.DB.
		Select("id, like_count, dislike_count").
		First(&comment, ID).Error; err != nil {
		return nil, ErrCommentNotFound
	}

	result := schemas.CommentStatsResponse{
		CommentID:    comment.ID,
		LikeCount:    comment.LikeCount,
		DislikeCount: comment.DislikeCount,
	}

	return &result, nil
}
