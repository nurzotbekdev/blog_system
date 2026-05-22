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

type CommentReactionController struct {
	CommentReactionService services.CommentReactionService
}

func NewCommentRectionController(reaction services.CommentReactionService) *CommentReactionController {
	return &CommentReactionController{CommentReactionService: reaction}
}

func (reaction *CommentReactionController) CreateCommentReaction(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized comment reaction attempt", zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	commentID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid comment id param", zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.CommentReactionRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid comment reaction JSON format", zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	logging.Log.Info("Creating comment reaction", zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID), zap.Bool("is_like", body.IsLike))
	newCommentRection := models.CommentReaction{
		UserID:    currentUserID,
		CommentID: commentID,
		IsLike:    body.IsLike,
	}

	if err := reaction.CommentReactionService.CreateCommentReaction(newCommentRection); err != nil {
		if errors.Is(err, services.ErrCommentNotFound) {
			logging.Log.Warn("Comment not found for reaction", zap.Uint("comment_id", commentID), zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrLikeAlredy) {
			logging.Log.Warn("Reaction already exists", zap.Uint("comment_id", commentID), zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Create comment reaction failed", zap.Uint("comment_id", commentID), zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Comment reaction created successfully", zap.Uint("comment_id", commentID), zap.Uint("user_id", currentUserID))
	ctx.JSON(http.StatusCreated, gin.H{"message": "Comment reaction created successfully"})
}

func (reaction *CommentReactionController) UpdateCommentReaction(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("unauthorized request to update comment reaction")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	commentReactionID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("iIvalid comment reaction id param", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.CommentReactionRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("invalid json format for comment reaction update", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Uint("reaction_id", commentReactionID))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := reaction.CommentReactionService.EditCommentReaction(currentUserID, commentReactionID, body.IsLike); err != nil {
		if errors.Is(err, services.ErrCommentReactionNotFound) {
			logging.Log.Warn("comment reaction not found", zap.Uint("user_id", currentUserID), zap.Uint("reaction_id", commentReactionID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("failed to update comment reaction", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Uint("reaction_id", commentReactionID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("comment reaction updated successfully", zap.Uint("user_id", currentUserID), zap.Uint("reaction_id", commentReactionID), zap.Bool("is_like", body.IsLike))
	ctx.JSON(http.StatusOK, gin.H{"message": "Comment reaction update successfully"})
}

func (reaction *CommentReactionController) DeleteCommentReaction(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized request to delete comment reaction")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	commentReactionID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid comment reaction id param", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := reaction.CommentReactionService.DeleteCommentReaction(currentUserID, commentReactionID); err != nil {
		if errors.Is(err, services.ErrCommentReactionNotFound) {
			logging.Log.Warn("Comment reaction not found", zap.Uint("user_id", currentUserID), zap.Uint("reaction_id", commentReactionID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrCommentNotFound) {
			logging.Log.Warn("Comment not found while deleting reaction", zap.Uint("user_id", currentUserID), zap.Uint("reaction_id", commentReactionID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to delete comment reaction", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Uint("reaction_id", commentReactionID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Comment reaction deleted successfully", zap.Uint("user_id", currentUserID), zap.Uint("reaction_id", commentReactionID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Comment reaction delete successfully"})
}

func (reaction *CommentReactionController) GetMyReaction(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized request to get my comment reaction")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	commentID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid comment id param", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reactionData, err := reaction.CommentReactionService.GetMyReaction(currentUserID, commentID)
	if err != nil {
		if errors.Is(err, services.ErrCommentReactionNotFound) {
			logging.Log.Warn("Comment reaction not found", zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to get my comment reaction", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Comment reaction fetched successfully", zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID))
	ctx.JSON(http.StatusOK, reactionData)
}
