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
		logging.Log.Warn("Unauthorized create video attempt")

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	currentUserID := userID.(uint)

	historyData, err := history.HistoryService.GetUserHistory(currentUserID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": historyData})
}

func (history *HistoryController) DeleteHistory(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized create video attempt")

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}
	currentUserID := userID.(uint)

	historyID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid subscription id param", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := history.HistoryService.DeleteHistory(currentUserID, historyID); err != nil {
		if errors.Is(err, services.ErrHistoryNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "History delete successfully"})
}
