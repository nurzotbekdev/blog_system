package services

import (
	"blog_system/config"
	"blog_system/logging"
	"blog_system/models"
	"blog_system/schemas"
	"errors"

	"go.uber.org/zap"
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
		logging.Log.Info("Creating comment reaction", zap.Uint("user_id", reaction.UserID), zap.Uint("comment_id", reaction.CommentID))
		var comment models.Comment
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&comment, reaction.CommentID).Error; err != nil {
			logging.Log.Warn("Comment not found for reaction", zap.Uint("comment_id", reaction.CommentID), zap.Error(err))
			return ErrCommentNotFound
		}

		var existing models.CommentReaction
		if err := tx.
			Where("user_id = ? AND comment_id = ?", reaction.UserID, reaction.CommentID).
			First(&existing).Error; err == nil {
			logging.Log.Warn("Reaction already exists", zap.Uint("user_id", reaction.UserID), zap.Uint("comment_id", reaction.CommentID))
			return ErrLikeAlready
		}

		if err := tx.Create(&reaction).Error; err != nil {
			logging.Log.Error("Failed to create comment reaction", zap.Uint("user_id", reaction.UserID), zap.Uint("comment_id", reaction.CommentID), zap.Error(err))
			return err
		}

		if reaction.IsLike {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("like_count",
					gorm.Expr("like_count + 1")).Error; err != nil {

				logging.Log.Error("Failed to increment like_count", zap.Uint("comment_id", reaction.CommentID), zap.Error(err))
				return err
			}
		} else {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("dislike_count",
					gorm.Expr("dislike_count + 1")).Error; err != nil {

				logging.Log.Error("Failed to increment dislike_count", zap.Uint("comment_id", reaction.CommentID), zap.Error(err))
				return err
			}
		}

		logging.Log.Info("Comment reaction created successfully", zap.Uint("user_id", reaction.UserID), zap.Uint("comment_id", reaction.CommentID))
		return nil
	})
}

func (s *commentReactionService) EditCommentReaction(userID, ID uint, isLike bool) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		logging.Log.Info("Editing comment reaction", zap.Uint("user_id", userID), zap.Uint("reaction_id", ID))

		var reaction models.CommentReaction
		if err := tx.Where("id = ? AND user_id = ?", ID, userID).First(&reaction).Error; err != nil {
			logging.Log.Warn("Comment reaction not found", zap.Uint("user_id", userID), zap.Uint("reaction_id", ID), zap.Error(err))
			return ErrCommentReactionNotFound
		}

		var comment models.Comment
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&comment, reaction.CommentID).Error; err != nil {
			logging.Log.Warn("Comment not found for reaction update", zap.Uint("comment_id", reaction.CommentID), zap.Error(err))
			return ErrCommentNotFound
		}

		if reaction.IsLike {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("like_count",
					gorm.Expr("GREATEST(like_count - 1,0)")).Error; err != nil {

				logging.Log.Error("Failed to decrement like_count", zap.Uint("comment_id", reaction.CommentID), zap.Error(err))
				return err
			}
		} else {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("dislike_count",
					gorm.Expr("GREATEST(dislike_count - 1,0)")).Error; err != nil {

				logging.Log.Error("Failed to decrement dislike_count", zap.Uint("comment_id", reaction.CommentID), zap.Error(err))
				return err
			}
		}

		if isLike {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("like_count",
					gorm.Expr("like_count + 1")).Error; err != nil {

				logging.Log.Error("Failed to increment like_count", zap.Uint("comment_id", reaction.CommentID), zap.Error(err))
				return err
			}
		} else {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("dislike_count",
					gorm.Expr("dislike_count + 1")).Error; err != nil {

				logging.Log.Error("Failed to increment dislike_count", zap.Uint("comment_id", reaction.CommentID), zap.Error(err))
				return err
			}
		}

		reaction.IsLike = isLike

		if err := tx.Save(&reaction).Error; err != nil {
			logging.Log.Error("Failed to save updated comment reaction", zap.Uint("reaction_id", ID), zap.Error(err))
			return err
		}

		logging.Log.Info("Comment reaction updated successfully", zap.Uint("user_id", userID), zap.Uint("reaction_id", ID))
		return nil
	})
}

func (s *commentReactionService) DeleteCommentReaction(userID, ID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		logging.Log.Info("Deleting comment reaction", zap.Uint("user_id", userID), zap.Uint("reaction_id", ID))

		var reaction models.CommentReaction
		if err := tx.Where("id = ? AND user_id = ?", ID, userID).First(&reaction).Error; err != nil {

			logging.Log.Warn("Comment reaction not found", zap.Uint("user_id", userID), zap.Uint("reaction_id", ID), zap.Error(err))
			return ErrCommentReactionNotFound
		}

		var comment models.Comment
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&comment, reaction.CommentID).Error; err != nil {

			logging.Log.Warn("Comment not found for reaction delete", zap.Uint("comment_id", reaction.CommentID), zap.Error(err))
			return ErrCommentNotFound
		}

		if err := tx.Delete(&reaction).Error; err != nil {

			logging.Log.Error("Failed to delete comment reaction", zap.Uint("reaction_id", ID), zap.Error(err))
			return err
		}

		if reaction.IsLike {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("like_count",
					gorm.Expr("GREATEST(like_count - 1,0)")).Error; err != nil {

				logging.Log.Error("Failed to decrement like_count", zap.Uint("comment_id", reaction.CommentID), zap.Error(err))
				return err
			}
		} else {
			if err := tx.Model(&models.Comment{}).
				Where("id = ?", reaction.CommentID).
				UpdateColumn("dislike_count",
					gorm.Expr("GREATEST(dislike_count - 1,0)")).Error; err != nil {

				logging.Log.Error("Failed to decrement dislike_count", zap.Uint("comment_id", reaction.CommentID), zap.Error(err))
				return err
			}
		}

		logging.Log.Info("Comment reaction deleted successfully", zap.Uint("user_id", userID), zap.Uint("reaction_id", ID))
		return nil
	})
}

func (s *commentReactionService) GetMyReaction(userID, ID uint) (*schemas.CommentReactionResponse, error) {
	logging.Log.Info("Fetching user comment reaction", zap.Uint("user_id", userID), zap.Uint("comment_id", ID))
	var reaction models.CommentReaction

	if err := config.DB.
		Where("user_id = ? AND comment_id = ?", userID, ID).
		First(&reaction).Error; err != nil {

		logging.Log.Warn("Comment reaction not found", zap.Uint("user_id", userID), zap.Uint("comment_id", ID), zap.Error(err))
		return nil, ErrCommentReactionNotFound
	}

	result := schemas.CommentReactionResponse{
		ID:        reaction.ID,
		IsLike:    reaction.IsLike,
		CreatedAt: reaction.CreatedAt,
	}

	logging.Log.Info("Comment reaction fetched successfully", zap.Uint("user_id", userID), zap.Uint("comment_id", ID))
	return &result, nil
}
