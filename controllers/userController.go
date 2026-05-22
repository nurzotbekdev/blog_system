package controllers

import (
	"blog_system/config"
	"blog_system/logging"
	"blog_system/models"
	"blog_system/security"
	"blog_system/services"
	"context"
	"encoding/json"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(user services.UserService) *UserController {
	return &UserController{UserService: user}
}

func (user *UserController) GoogleLogin(ctx *gin.Context) {
	logging.Log.Info("Google login initiated", zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
	tokenStr, err := ctx.Cookie("token")
	if err == nil {
		logging.Log.Info("Token found in cookie, validating token")
		_, err := security.ValidateToken(tokenStr)
		if err == nil {
			logging.Log.Info("User already authenticated, redirecting to frontend")
			ctx.Redirect(302, "http://localhost:3000")
			return
		}

		logging.Log.Warn("Invalid token in cookie", zap.Error(err))
	}

	url := config.GoogleOAuthConfig.AuthCodeURL("state-token")

	logging.Log.Info("Redirecting to Google OAuth", zap.String("redirect_url", url))
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (user *UserController) GoogleCallback(ctx *gin.Context) {
	logging.Log.Info("Google callback started", zap.String("path", ctx.FullPath()))
	code := ctx.Query("code")
	if code == "" {
		logging.Log.Warn("Google callback missing code")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "code not found"})
		return
	}

	logging.Log.Info("Exchanging OAuth code for token")
	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		logging.Log.Error("OAuth token exchange failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "token exchange failed"})
		return
	}

	client := config.GoogleOAuthConfig.Client(context.Background(), token)
	logging.Log.Info("Fetching Google user info")

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		logging.Log.Error("Failed to get Google user info", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		logging.Log.Error("Failed to decode Google user info", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode user info"})
		return
	}

	logging.Log.Info("Google user received", zap.String("email", googleUser.Email), zap.String("google_id", googleUser.ID))
	newUser := models.User{
		GoogleID:     googleUser.ID,
		Email:        googleUser.Email,
		FullName:     googleUser.Name,
		ProfileImage: googleUser.Picture,
		Role:         "user",
	}

	savedUser, err := user.UserService.SignIn(newUser)
	if err != nil {
		logging.Log.Error("User save/sign-in failed", zap.String("email", googleUser.Email), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user"})
		return
	}

	logging.Log.Info("User authenticated successfully", zap.Uint("user_id", savedUser.ID), zap.String("email", savedUser.Email))
	jwtToken, _ := security.GenerateToken(savedUser.ID)
	ctx.SetCookie("token", jwtToken, 3600*24, "/", "", true, true)

	logging.Log.Info("User logged in via Google successfully", zap.Uint("user_id", savedUser.ID))
	ctx.Redirect(302, "http://localhost:3000")
}

func (user *UserController) MyProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		logging.Log.Warn("Unauthorized profile access attempt", zap.String("path", ctx.FullPath()))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	logging.Log.Info("Fetching user profile", zap.Any("user_id", userID))
	var currentUser models.User
	err := config.DB.First(&currentUser, userID).Error
	if err != nil {
		logging.Log.Warn("User not found in DB", zap.Any("user_id", userID), zap.Error(err))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found in DB"})
		return
	}

	logging.Log.Info("User profile fetched successfully", zap.Uint("user_id", currentUser.ID), zap.String("email", currentUser.Email))
	ctx.JSON(http.StatusOK, gin.H{
		"id":            currentUser.ID,
		"email":         currentUser.Email,
		"name":          currentUser.FullName,
		"profile_image": currentUser.ProfileImage,
		"created_at":    currentUser.CreatedAt,
	})
}
