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

type LanguageController struct {
	LanguageService services.LanguageService
}

func NewLanguageController(language services.LanguageService) *LanguageController {
	return &LanguageController{LanguageService: language}
}

func (language *LanguageController) Create(ctx *gin.Context) {
	var body schemas.LanguageSchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid json format for language create", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	newLanguage := models.Language{
		Name: body.Name,
		Code: body.Code,
	}

	if err := language.LanguageService.CreateLanguage(newLanguage); err != nil {
		if errors.Is(err, services.ErrAlreadyLanguage) {
			logging.Log.Warn("Duplicate language detected", zap.String("name", body.Name), zap.String("code", body.Code))
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to create language", zap.Error(err), zap.String("name", body.Name), zap.String("code", body.Code))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Language created successfully", zap.String("name", body.Name), zap.String("code", body.Code))
	ctx.JSON(http.StatusCreated, gin.H{"message": "Language created"})
}

func (language *LanguageController) AllLanguage(ctx *gin.Context) {
	languageData, err := language.LanguageService.GetAllLanguage()
	if err != nil {
		logging.Log.Error("Failed to fetch languages", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Languages fetched successfully", zap.Int("count", len(languageData)), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
	ctx.JSON(http.StatusOK, gin.H{"data": languageData})
}

func (language *LanguageController) UpdateLanguage(ctx *gin.Context) {
	languageID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid language id param", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.LanguageSchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid json format for language update", zap.Error(err), zap.Uint("language_id", languageID))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := language.LanguageService.EditLanguage(languageID, body.Name, body.Code); err != nil {
		if errors.Is(err, services.ErrAlreadyLanguage) {
			logging.Log.Warn("Duplicate language detected", zap.Uint("language_id", languageID), zap.String("name", body.Name), zap.String("code", body.Code))
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrLanguageNotFound) {
			logging.Log.Warn("Language not found", zap.Uint("language_id", languageID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to update language", zap.Error(err), zap.Uint("language_id", languageID), zap.String("name", body.Name), zap.String("code", body.Code))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Language updated successfully", zap.Uint("language_id", languageID), zap.String("name", body.Name), zap.String("code", body.Code))
	ctx.JSON(http.StatusOK, gin.H{"message": "Language update successful"})
}

func (language *LanguageController) RemoveLanguage(ctx *gin.Context) {
	languageID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid language id param", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := language.LanguageService.DeleteLanguage(languageID); err != nil {
		if errors.Is(err, services.ErrLanguageNotFound) {
			logging.Log.Warn("Language not found", zap.Uint("language_id", languageID))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Failed to delete language", zap.Error(err), zap.Uint("language_id", languageID))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Language deleted successfully", zap.Uint("language_id", languageID))
	ctx.JSON(http.StatusOK, gin.H{"message": "Language delete successful"})
}
