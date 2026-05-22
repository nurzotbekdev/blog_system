package controllers

import (
	"blog_system/helper"
	"blog_system/logging"
	"blog_system/models"
	"blog_system/schemas"
	"blog_system/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommentController struct {
	CommentService services.CommentService
}

func NewCommentController(comment services.CommentService) *CommentController {
	return &CommentController{CommentService: comment}
}

func (comment *CommentController) Create(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized comment create attempt", zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	var body schemas.CommentSchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid comment JSON format", zap.Uint("user_id", currentUserID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	logging.Log.Info("Creating comment", zap.Uint("user_id", currentUserID), zap.Uint("video_id", body.VideoID))
	newComment := models.Comment{
		UserID:   currentUserID,
		VideoID:  body.VideoID,
		Content:  body.Content,
		ParentID: body.ParentID,
	}

	if err := comment.CommentService.CreateComment(newComment); err != nil {
		if errors.Is(err, services.ErrVideoNotFound) {
			logging.Log.Warn("Video not found for comment", zap.Uint("video_id", body.VideoID), zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrParentCommentNotFound) {
			logging.Log.Warn("Parent comment not found", zap.Uint("user_id", currentUserID), zap.Any("parent_id", body.ParentID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Create comment failed", zap.Uint("user_id", currentUserID), zap.Uint("video_id", body.VideoID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Comment created successfully", zap.Uint("user_id", currentUserID), zap.Uint("video_id", body.VideoID))
	ctx.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully"})
}

func (comment *CommentController) GetVideoComments(ctx *gin.Context) {
	videoID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid video id param for comments", zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	sort := ctx.Query("sort")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	logging.Log.Info("Fetching video comments", zap.Uint("video_id", videoID), zap.Int("page", page), zap.Int("limit", limit), zap.String("sort", sort))
	commentData, err := comment.CommentService.GetVideoComments(page, limit, videoID, sort)
	if err != nil {
		logging.Log.Error("Failed to fetch video comments", zap.Uint("video_id", videoID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Video comments fetched successfully", zap.Uint("video_id", videoID), zap.Int("page", page))
	ctx.JSON(http.StatusOK, gin.H{"data": commentData})
}

func (comment *CommentController) UpdateComment(ctx *gin.Context) {
	commentID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid comment id param", zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized update comment attempt", zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	var body schemas.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid comment update JSON format", zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	logging.Log.Info("Updating comment", zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID))
	if err := comment.CommentService.EditVideoComment(currentUserID, commentID, body.Content); err != nil {
		if errors.Is(err, services.ErrCommentNotFound) {
			logging.Log.Warn("Comment not found for update", zap.Uint("comment_id", commentID), zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Update comment failed", zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Comment updated successfully", zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Comment update successfully"})
}

func (comment *CommentController) DeleteComment(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized delete comment attempt", zap.String("path", ctx.FullPath()))
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

	logging.Log.Info("Deleting comment", zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID))
	if err := comment.CommentService.DeleteComment(currentUserID, commentID); err != nil {
		if errors.Is(err, services.ErrCommentNotFound) {
			logging.Log.Warn("Comment not found", zap.Uint("comment_id", commentID), zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrVideoNotFound) {
			logging.Log.Warn("Video not found for comment delete", zap.Uint("comment_id", commentID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrChannelNotFound) {
			logging.Log.Warn("Channel not found for comment delete", zap.Uint("comment_id", commentID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrForbidden) {
			logging.Log.Warn("Forbidden comment delete attempt", zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID))
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Delete comment failed", zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Comment deleted successfully", zap.Uint("user_id", currentUserID), zap.Uint("comment_id", commentID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Comment delete successfully"})
}

func (comment *CommentController) ReplyComment(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized reply comment attempt", zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	parentID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid parent comment id param", zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid reply comment JSON format", zap.Uint("user_id", currentUserID), zap.Uint("parent_comment_id", parentID), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	logging.Log.Info("Creating reply comment", zap.Uint("user_id", currentUserID), zap.Uint("parent_comment_id", parentID))
	newComment := models.Comment{
		UserID:  currentUserID,
		Content: body.Content,
	}

	if err := comment.CommentService.CreateReply(parentID, newComment); err != nil {
		if errors.Is(err, services.ErrParentCommentNotFound) {
			logging.Log.Warn("Parent comment not found for reply", zap.Uint("parent_comment_id", parentID), zap.Uint("user_id", currentUserID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Create reply comment failed", zap.Uint("user_id", currentUserID), zap.Uint("parent_comment_id", parentID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Reply comment created successfully", zap.Uint("user_id", currentUserID), zap.Uint("parent_comment_id", parentID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Comment reply successfully"})
}

func (comment *CommentController) GetStats(ctx *gin.Context) {
	commentID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid comment id param", zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logging.Log.Info("Fetching comment stats", zap.Uint("comment_id", commentID))
	commentStats, err := comment.CommentService.GetStats(commentID)
	if err != nil {
		if errors.Is(err, services.ErrCommentNotFound) {
			logging.Log.Warn("Comment not found for stats", zap.Uint("comment_id", commentID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to get comment stats", zap.Uint("comment_id", commentID), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Comment stats fetched successfully", zap.Uint("comment_id", commentID))
	ctx.JSON(http.StatusOK, commentStats)
}
