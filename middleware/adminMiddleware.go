package middleware

import (
	"blog_system/logging"
	"blog_system/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInterface, exists := ctx.Get("user")
		if !exists {
			logging.Log.Warn("Unauthorized access: user not found")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			ctx.Abort()
			return
		}

		user, ok := userInterface.(models.User)
		if !ok {
			logging.Log.Error("Invalid user type in context", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			ctx.Abort()
			return
		}

		if user.Role != "admin" {
			logging.Log.Warn("Forbidden: non-admin access attempt", zap.Uint("user_id", user.ID), zap.String("role", user.Role))
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have administrator rights"})
			ctx.Abort()
			return
		}

		logging.Log.Info("Admin access granted", zap.Uint("user_id", user.ID), zap.String("path", ctx.FullPath()))
		ctx.Next()
	}
}
