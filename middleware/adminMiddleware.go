package middleware

import (
	"blog_system/config"

	"blog_system/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDRaw, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			ctx.Abort()
			return
		}

		userID, ok := userIDRaw.(uint)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			ctx.Abort()
			return
		}

		var user models.User
		if err := config.DB.First(&user, userID).Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			ctx.Abort()
			return
		}

		if user.Role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have administrator rights"})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
