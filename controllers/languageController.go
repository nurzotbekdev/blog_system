package controllers

import (
	"blog_system/helper"
	"blog_system/models"
	"blog_system/schemas"
	"blog_system/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	newLanguage := models.Language{
		Name: body.Name,
		Code: body.Code,
	}

	if err := language.LanguageService.CreateLanguage(newLanguage); err != nil {
		if errors.Is(err, services.ErrAlreadyLanguage) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Language created"})
}

func (language *LanguageController) AllLanguage(ctx *gin.Context) {
	languageData, err := language.LanguageService.GetAllLanguage()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": languageData})
}

func (language *LanguageController) UpdateLanguage(ctx *gin.Context) {
	languageID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.LanguageSchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := language.LanguageService.EditLanguage(languageID, body.Name, body.Code); err != nil {
		if errors.Is(err, services.ErrAlreadyLanguage) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrLanguageNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Language update successful"})
}

func (language *LanguageController) RemoveLanguage(ctx *gin.Context) {
	languageID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := language.LanguageService.DeleteLanguage(languageID); err != nil {
		if errors.Is(err, services.ErrLanguageNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Language delete successful"})
}
