package middleware

import (
	"blog_system/config"
	"blog_system/logging"
	"blog_system/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDRaw, exists := ctx.Get("user_id")
		if !exists {
			logging.Log.Warn("Unauthorized access: user not found")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			ctx.Abort()
			return
		}

		userID, ok := userIDRaw.(uint)
		if !ok {
			logging.Log.Error("Invalid user_id type")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			ctx.Abort()
			return
		}

		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			logging.Log.Warn("User not found in DB", zap.Uint("user_id", userID))
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			ctx.Abort()
			return
		}

		if user.Role != "admin" {
			logging.Log.Warn("Forbidden access", zap.Uint("user_id", user.ID), zap.String("role", user.Role))
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have administrator rights"})
			ctx.Abort()
			return
		}

		logging.Log.Info("Admin access granted", zap.Uint("user_id", user.ID))
		ctx.Set("user", user)
		ctx.Next()
	}
}
