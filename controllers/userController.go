package controllers

import (
	"blog_system/config"
	"blog_system/models"
	"blog_system/security"
	"blog_system/services"
	"context"
	"encoding/json"

	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(user services.UserService) *UserController {
	return &UserController{UserService: user}
}

func (user *UserController) GoogleLogin(ctx *gin.Context) {
	tokenStr, err := ctx.Cookie("token")

	if err == nil {
		_, err := security.ValidateToken(tokenStr)
		if err == nil {
			ctx.Redirect(302, "http://localhost:3000")
			return
		}
	}

	url := config.GoogleOAuthConfig.AuthCodeURL("state-token")
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (user *UserController) GoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "code not found",
		})
		return
	}

	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "token exchange failed",
		})
		return
	}
	client := config.GoogleOAuthConfig.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user info",
		})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to decode user info",
		})
		return
	}

	newUser := models.User{
		GoogleID:     googleUser.ID,
		Email:        googleUser.Email,
		FullName:     googleUser.Name,
		ProfileImage: googleUser.Picture,
		Role:         "user",
	}

	savedUser, err := user.UserService.SignIn(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save user",
		})
		return
	}

	jwtToken, _ := security.GenerateToken(savedUser.ID)
	ctx.SetCookie("token", jwtToken, 3600*24, "/", "", true, true)

	ctx.Redirect(302, "http://localhost:3000")
}

func (user *UserController) MyProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	var currentUser models.User
	err := config.DB.First(&currentUser, userID).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found in DB",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":            currentUser.ID,
		"email":         currentUser.Email,
		"name":          currentUser.FullName,
		"profile_image": currentUser.ProfileImage,
		"created_at":    currentUser.CreatedAt,
	})
}
