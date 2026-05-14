package services

import (
	"blog_system/config"
	"blog_system/models"
	"blog_system/schemas"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentReactionService interface {
	CreateCommentReaction(reaction models.CommentReaction) error
	EditCommentReaction(userID, ID uint, isLike bool) error
	DeleteCommentReaction(userID, ID uint) error
	GetMyReaction(userID, ID uint) (*schemas.CommentReactionResponse, error)
}

type commentReactionService struct{}

func NewCommentRectionServices() CommentReactionService {
	return &commentReactionService{}
}

var (
	ErrLikeAlredy              = errors.New("Like already")
	ErrCommentReactionNotFound = errors.New("Comment reaction not found")
)

func (s *commentReactionService) CreateCommentReaction(reaction models.CommentReaction) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var comment models.Comment
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&comment, reaction.CommentID).Error; err != nil {
			return ErrCommentNotFound
		}

		var existing models.CommentReaction
		err := tx.
			Where("user_id = ? AND comment_id = ?", reaction.UserID, reaction.CommentID).
			First(&existing).Error

		if err == nil {

			if existing.IsLike == reaction.IsLike {
				return ErrLikeAlredy
			}

			if existing.IsLike {
				if err := tx.Model(&models.Comment{}).
					Where("id = ?", reaction.CommentID).
					UpdateColumn("like_count",
						gorm.Expr("GREATEST(like_count - 1,0)")).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(&models.Comment{}).
					Where("id = ?", reaction.CommentID).
					UpdateColumn("dislike_count",
						gorm.Expr("GREATEST(dislike_count - 1,0)")).Error; err != nil {
					return err
				}
			}

			existing.IsLike = reaction.IsLike
			return tx.Save(&existing).Error
		}

		if err := tx.Create(&reaction).Error; err != nil {
			return err
		}

		if reaction.IsLike {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("like_count",
					gorm.Expr("like_count + 1")).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("dislike_count",
					gorm.Expr("dislike_count + 1")).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *commentReactionService) EditCommentReaction(userID, ID uint, isLike bool) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var reaction models.CommentReaction
		if err := tx.Where("id = ? AND user_id = ?", ID, userID).First(&reaction).Error; err != nil {
			return ErrCommentReactionNotFound
		}

		var comment models.Comment
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&comment, reaction.CommentID).Error; err != nil {
			return ErrCommentNotFound
		}

		if reaction.IsLike {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("like_count",
					gorm.Expr("GREATEST(like_count - 1,0)")).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("dislike_count",
					gorm.Expr("GREATEST(dislike_count - 1,0)")).Error; err != nil {
				return err
			}
		}

		if isLike {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("like_count",
					gorm.Expr("like_count + 1")).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("dislike_count",
					gorm.Expr("dislike_count + 1")).Error; err != nil {
				return err
			}
		}

		reaction.IsLike = isLike
		return tx.Save(&reaction).Error
	})
}

func (s *commentReactionService) DeleteCommentReaction(userID, ID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var reaction models.CommentReaction
		if err := tx.Where("id = ? AND user_id = ?", ID, userID).First(&reaction).Error; err != nil {
			return ErrCommentReactionNotFound
		}

		var comment models.Comment
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&comment, reaction.CommentID).Error; err != nil {
			return ErrCommentNotFound
		}

		if err := tx.Delete(&reaction).Error; err != nil {
			return err
		}

		if reaction.IsLike {
			tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("like_count",
					gorm.Expr("GREATEST(like_count - 1,0)"))
		} else {
			tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("dislike_count",
					gorm.Expr("GREATEST(dislike_count - 1,0)"))
		}

		return nil
	})
}

func (s *commentReactionService) GetMyReaction(userID, ID uint) (*schemas.CommentReactionResponse, error) {
	var reaction models.CommentReaction

	if err := config.DB.
		Where("user_id = ? AND comment_id = ?", userID, ID).
		First(&reaction).Error; err != nil {

		return nil, ErrCommentReactionNotFound
	}

	result := schemas.CommentReactionResponse{
		ID:        reaction.ID,
		IsLike:    reaction.IsLike,
		CreatedAt: reaction.CreatedAt,
	}

	return &result, nil
}
