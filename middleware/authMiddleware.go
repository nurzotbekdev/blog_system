package middleware

import (
	"blog_system/logging"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("token")
		if err != nil {
			logging.Log.Warn("Token not found", zap.String("ip", ctx.ClientIP()))
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil || !token.Valid {
			logging.Log.Warn("Invalid token", zap.Error(err))
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalid"})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			ctx.Abort()
			return
		}

		sub, ok := claims["sub"].(float64)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			ctx.Abort()
			return
		}

		userID := uint(sub)

		ctx.Set("user_id", userID)

		logging.Log.Info("Authenticated request",
			zap.Uint("user_id", userID),
		)

		ctx.Next()
	}
}
