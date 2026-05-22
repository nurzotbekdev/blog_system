package controllers

import (
	"blog_system/helper"
	"blog_system/logging"
	"blog_system/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HistoryController struct {
	HistoryService services.HistoryService
}

func NewHistoryController(history services.HistoryService) *HistoryController {
	return &HistoryController{HistoryService: history}
}

func (history *HistoryController) GetUserHistory(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized request to get user history", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	historyData, err := history.HistoryService.GetUserHistory(currentUserID, page, limit)
	if err != nil {
		logging.Log.Error("Failed to fetch user history", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Int("page", page), zap.Int("limit", limit))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("User history fetched successfully", zap.Uint("user_id", currentUserID), zap.Int("page", page), zap.Int("limit", limit), zap.Int("count", len(historyData)))
	ctx.JSON(http.StatusOK, gin.H{"data": historyData})
}

func (history *HistoryController) DeleteHistory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized request to delete history", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(uint)

	historyID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("invalid history id param", zap.Error(err), zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := history.HistoryService.DeleteHistory(currentUserID, historyID); err != nil {
		if errors.Is(err, services.ErrHistoryNotFound) {
			logging.Log.Warn("history not found", zap.Uint("user_id", currentUserID), zap.Uint("history_id", historyID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("failed to delete history", zap.Error(err), zap.Uint("user_id", currentUserID), zap.Uint("history_id", historyID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("History deleted successfully", zap.Uint("user_id", currentUserID), zap.Uint("history_id", historyID))
	ctx.JSON(http.StatusOK, gin.H{"message": "History delete successfully"})
}
