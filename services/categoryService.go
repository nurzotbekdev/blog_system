package services

import (
	"blog_system/config"
	"blog_system/logging"
	"blog_system/models"
	"blog_system/schemas"
	"errors"

	"go.uber.org/zap"
)

type CategoryService interface {
	CreateCategory(category models.Category) error
	GetAllCategory() ([]schemas.CategoryResponse, error)
	EditCategory(ID uint, name string) error
	DeleteCategory(ID uint) error
}

type categoryService struct{}

func NewCategoryServices() CategoryService {
	return &categoryService{}
}

var (
	ErrAlreadyCategory  = errors.New("Category already added database")
	ErrCategoryNotFound = errors.New("Category not found")
)

func (s *categoryService) CreateCategory(category models.Category) error {
	var existing models.Category
	if err := config.DB.Where("name= ?", category.Name).First(&existing).Error; err == nil {
		logging.Log.Warn("Category already added database", zap.Uint("category_id", category.ID), zap.Error(err))
		return ErrAlreadyCategory
	}

	if err := config.DB.Create(&category).Error; err != nil {
		logging.Log.Error("Failed to create category", zap.Uint("category_id", category.ID), zap.Error(err))
		return err
	}

	logging.Log.Info("Category created successfully", zap.Uint("category_id", category.ID))
	return nil
}

func (s *categoryService) GetAllCategory() ([]schemas.CategoryResponse, error) {
	var categorys []models.Category
	if err := config.DB.Find(&categorys).Error; err != nil {
		logging.Log.Error("Failed to get categories", zap.Error(err))
		return nil, err
	}

	var response []schemas.CategoryResponse
	for _, c := range categorys {
		response = append(response, schemas.CategoryResponse{
			ID:        c.ID,
			Name:      c.Name,
			CreatedAt: c.CreatedAt,
		})
	}

	logging.Log.Info("Categories fetched successfully", zap.Int("count", len(response)))
	return response, nil
}

func (s *categoryService) EditCategory(ID uint, name string) error {
	var category models.Category

	if err := config.DB.Where("id = ?", ID).First(&category).Error; err != nil {
		logging.Log.Warn("Category not found", zap.Uint("category_id", ID), zap.Error(err))
		return ErrCategoryNotFound
	}

	var existing models.Category
	if err := config.DB.
		Where("name = ? AND id != ?", name, ID).
		First(&existing).Error; err == nil {

		logging.Log.Warn("Category name already exists", zap.String("name", name))
		return ErrAlreadyCategory
	}

	if err := config.DB.Model(&category).Update("name", name).Error; err != nil {
		logging.Log.Error("Failed to update category", zap.Uint("category_id", ID), zap.Error(err))
		return err
	}

	logging.Log.Info("Category updated successfully", zap.Uint("category_id", ID))
	return nil
}

func (s *categoryService) DeleteCategory(ID uint) error {
	var category models.Category
	if err := config.DB.Where("id = ?", ID).First(&category).Error; err != nil {
		logging.Log.Warn("Category not found", zap.Uint("category_id", ID), zap.Error(err))
		return ErrCategoryNotFound
	}

	if err := config.DB.Unscoped().Delete(&category).Error; err != nil {
		logging.Log.Error("Failed to delete category", zap.Uint("category_id", category.ID), zap.Error(err))
		return err
	}

	logging.Log.Info("Category deleted successfully", zap.Uint("category_id", category.ID))
	return nil
}
