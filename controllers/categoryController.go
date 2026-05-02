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

type CategoryController struct {
	CategoryService services.CategoryService
}

func NewCategoryController(category services.CategoryService) *CategoryController {
	return &CategoryController{CategoryService: category}
}

func (category *CategoryController) Create(ctx *gin.Context) {
	var body schemas.CategorySchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid JSON format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	newCategory := models.Category{
		Name: body.Name,
	}

	if err := category.CategoryService.CreateCategory(newCategory); err != nil {
		if errors.Is(err, services.ErrAlreadyCategory) {
			logging.Log.Warn("Duplicate subscription", zap.Error(err))
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Update category failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Category created successfully")
	ctx.JSON(http.StatusCreated, gin.H{"message": "Category created"})
}

func (category *CategoryController) AllCategory(ctx *gin.Context) {
	categoryData, err := category.CategoryService.GetAllCategory()

	if err != nil {
		logging.Log.Error("Failed to fetch categories", zap.Error(err), zap.String("path", ctx.FullPath()), zap.String("method", ctx.Request.Method))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Categories fetched successfully", zap.Int("count", len(categoryData)), zap.String("path", ctx.FullPath()))
	ctx.JSON(http.StatusOK, gin.H{"data": categoryData})
}

func (category *CategoryController) UpdateCategory(ctx *gin.Context) {
	categoryID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid subscription id param", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body schemas.CategorySchemas
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logging.Log.Warn("Invalid JSON format", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := category.CategoryService.EditCategory(categoryID, body.Name); err != nil {
		if errors.Is(err, services.ErrAlreadyCategory) {
			logging.Log.Warn("Duplicate category", zap.Error(err))
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, services.ErrCategoryNotFound) {
			logging.Log.Warn("Category not found", zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Update category failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Category update successfully")
	ctx.JSON(http.StatusOK, gin.H{"message": "Category update successful"})
}

func (category *CategoryController) RemoveCategory(ctx *gin.Context) {
	categoryID, err := helper.ParseUintParam(ctx, "id")
	if err != nil {
		logging.Log.Warn("Invalid subscription id param", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := category.CategoryService.DeleteCategory(categoryID); err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			logging.Log.Warn("Category not found", zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logging.Log.Error("Update category failed", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	logging.Log.Info("Category delete successfully")
	ctx.JSON(http.StatusOK, gin.H{"message": "Category delete successful"})
}
